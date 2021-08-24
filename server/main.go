package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Name     string  `gorm:"unique;not null;type:varchar(20)"`
	Email    string  `gorm:"unique;not null;type:varchar(100)"`
	Password []byte  `gorm:"not null"`
	Tweets   []Tweet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Tweet struct {
	gorm.Model
	Content string `gorm:"not null;type:text"`
	UserId  uint   `gorm:"not null;"`
	User    User
}

type Follow struct {
	gorm.Model
	FollowerId uint `gorm:"not null"`
	Follower   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FollowedId uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Followed   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func dbConnect() *gorm.DB {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed DB connect")
	}
	return db
}

func createUser(name string, email string, password string, db *gorm.DB) (User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := User{Name: name, Email: email, Password: hashedPassword}
	result := db.Create(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func loginUser(loginUserId uint, c *gin.Context, db *gorm.DB) (User, *gorm.DB) {
	var user User
	err := db.First(&user, loginUserId)
	if err.Error != nil {
		return user, err
	}
	return user, nil
}

func authorize(email string, password string, c *gin.Context, db *gorm.DB) (User, error) {
	var user User
	db_err := db.Where("email = ?", email).First(&user)
	if db_err.Error != nil {
		return user, db_err.Error
	}
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return user, err
	}
	session := sessions.Default(c)
	session.Set("UserId", user.ID)
	session.Save()
	return user, nil
}

func createTweet(content string, userId uint, db *gorm.DB) (Tweet, error) {
	tweet := Tweet{Content: content, UserId: userId}
	result := db.Create(&tweet)

	if result.Error != nil {
		return tweet, result.Error
	}
	return tweet, nil
}

func getUsers(loginUserId uint, db *gorm.DB) ([]User, error) {
	var users []User
	var follows []Follow
	var followUserIds []uint

	db.Select("followed_id").Where("follower_id = ?", loginUserId).Find(&follows).Pluck("followed_id", &followUserIds)
	// TODO: impl pagination
	result := db.Where("id NOT IN (?)", append(followUserIds, loginUserId)).Order("id DESC").Find(&users)

	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func createFollow(followerId uint, followedId uint, db *gorm.DB) (Follow, error) {
	follow := Follow{FollowerId: followerId, FollowedId: followedId}
	result := db.Create(&follow)

	if result.Error != nil {
		return follow, result.Error
	}
	return follow, nil
}

func getTweets(loginUserId uint, db *gorm.DB) ([]Tweet, error) {
	var tweets []Tweet
	var follows []Follow
	var followUserIds []uint

	db.Select("followed_id").Where("follower_id = ?", loginUserId).Find(&follows).Pluck("followed_id", &followUserIds)
	// TODO: impl pagination
	result := db.Preload("User").Where("user_id IN (?)", append(followUserIds, loginUserId)).Order("id DESC").Find(&tweets)

	if result.Error != nil {
		return tweets, result.Error
	}
	return tweets, nil
}

func getSessionUserId(c *gin.Context) interface{} {
	session := sessions.Default(c)
	return session.Get("UserId")
}

func loginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := getSessionUserId(c)

		if UserId == nil {
			c.HTML(http.StatusOK, "signin.html", gin.H{
				"err": "ログインしてください",
			})
			return
		}
		c.Next()
	}
}

func main() {
	db := dbConnect()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Tweet{}, &Follow{})

	// routing
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	// session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	loginGroup := router.Group("/")
	loginGroup.Use(loginCheckMiddleware())
	{
		loginGroup.POST("/tweet", func(c *gin.Context) {
			UserId := getSessionUserId(c)
			content := c.PostForm("content")
			_, err := createTweet(content, UserId.(uint), db)
			if err != nil {
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"err": err.Error()})
				return
			}
			c.Redirect(http.StatusFound, "/")
		})
		loginGroup.GET("/users", func(c *gin.Context) {
			UserId := getSessionUserId(c)
			users, err := getUsers(UserId.(uint), db)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "users.html", gin.H{"err": err.Error()})
				return
			}
			c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
		})

		loginGroup.POST("/follow", func(c *gin.Context) {
			UserId := getSessionUserId(c)
			FollowedId, _ := strconv.ParseUint(c.PostForm("userId"), 10, 64)
			_, err := createFollow(UserId.(uint), uint(FollowedId), db)
			if err != nil {
				c.HTML(http.StatusBadRequest, "users.html", gin.H{"err": err.Error()})
				return
			}
			c.Redirect(http.StatusFound, "/users")
		})

	}
	router.GET("/", func(c *gin.Context) {
		UserId := getSessionUserId(c)
		if UserId == nil {
			c.HTML(http.StatusOK, "signin.html", gin.H{
				"err": "ログインしてください",
			})
			return
		}
		user, err := loginUser(UserId.(uint), c, db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "signin.html", gin.H{
				"err": err,
			})
			return
		}
		tweets, error := getTweets(UserId.(uint), db)
		if error != nil {
			c.HTML(http.StatusInternalServerError, "signin.html", gin.H{
				"err": error,
			})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name":   user.Name,
			"tweets": tweets,
		})
	})
	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})
	router.POST("/signup", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")
		user, err := createUser(name, email, password, db)
		if err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err.Error()})
			return
		}
		session := sessions.Default(c)
		session.Set("UserId", user.ID)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})
	router.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", gin.H{})
	})
	router.POST("/signin", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		_, err := authorize(email, password, c, db)
		if err != nil {
			c.HTML(http.StatusBadRequest, "signin.html", gin.H{"err": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/")
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
