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
	Name     string `gorm:"unique;not null;type:varchar(20)"`
	Email    string `gorm:"unique;not null;type:varchar(100)"`
	Password []byte `gorm:"not null"`
}

type Tweet struct {
	gorm.Model
	Content string `gorm:"not null;type:text"`
	UserId  uint   `gorm:"not null;constrant:OnDelete:CASCADE"`
	User    User
}

type Follow struct {
	gorm.Model
	FollowerId uint `gorm:"not null;constrant:OnDelete:CASCADE"`
	Follower   User
	FollowedId uint `gorm:"not null;constrant:OnDelete:CASCADE"`
	Followed   User
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

func createUser(name string, email string, password string) (User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := User{Name: name, Email: email, Password: hashedPassword}
	db := dbConnect()
	result := db.Create(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func loginUser(c *gin.Context) (User, *gorm.DB) {
	db := dbConnect()
	session := sessions.Default(c)
	UserId := session.Get("UserId")
	var user User
	err := db.First(&user, UserId)
	if err.Error != nil {
		return user, err
	}
	return user, nil
}

func authorize(email string, password string, c *gin.Context) (User, error) {
	db := dbConnect()
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

func createTweet(content string, userId uint) (Tweet, error) {
	db := dbConnect()
	tweet := Tweet{Content: content, UserId: userId}
	result := db.Create(&tweet)

	if result.Error != nil {
		return tweet, result.Error
	}
	return tweet, nil
}

func getUsers(loginUserId uint) ([]User, error) {
	db := dbConnect()
	var users []User
	var follows []Follow
	// TODO: impl pagination
	followUsers := db.Select("followed_id").Where("follower_id = ?", loginUserId).Find(&follows)
	result := db.Where("id NOT IN (?)", followUsers).Not("id = ?", loginUserId).Order("id DESC").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func createFollow(followerId uint, followedId uint) (Follow, error) {
	db := dbConnect()
	follow := Follow{FollowerId: followerId, FollowedId: followedId}
	result := db.Create(&follow)

	if result.Error != nil {
		return follow, result.Error
	}
	return follow, nil
}

func getTweets(loginUserId uint) ([]Tweet, error) {
	db := dbConnect()
	var tweets []Tweet
	var follows []Follow
	// TODO: impl pagination
	followUsers := db.Select("followed_id").Where("follower_id = ?", loginUserId).Find(&follows)
	result := db.Preload("User").Where("user_id IN (?)", followUsers).Or("user_id = ?", loginUserId).Order("id DESC").Find(&tweets)
	if result.Error != nil {
		return tweets, result.Error
	}
	return tweets, nil
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

	router.GET("/", func(c *gin.Context) {
		user, err := loginUser(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "signin.html", gin.H{
				"err": err,
			})
			return
		}
		tweets, error := getTweets(user.ID)
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
		user, err := createUser(name, email, password)
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
		_, err := authorize(email, password, c)
		if err != nil {
			c.HTML(http.StatusBadRequest, "signin.html", gin.H{"err": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/")
	})
	router.POST("/tweet", func(c *gin.Context) {
		session := sessions.Default(c)
		UserId := session.Get("UserId").(uint)
		content := c.PostForm("content")
		_, err := createTweet(content, UserId)
		if err != nil {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"err": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/")
	})
	router.GET("/users", func(c *gin.Context) {
		session := sessions.Default(c)
		UserId := session.Get("UserId").(uint)
		users, err := getUsers(UserId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "users.html", gin.H{"err": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
	})

	router.POST("/follow", func(c *gin.Context) {
		session := sessions.Default(c)
		LoginUserId := session.Get("UserId").(uint)
		UserId, _ := strconv.ParseUint(c.PostForm("userId"), 10, 64)
		_, err := createFollow(LoginUserId, uint(UserId))
		if err != nil {
			c.HTML(http.StatusBadRequest, "users.html", gin.H{"err": err.Error()})
			return
		}
		c.Redirect(http.StatusFound, "/users")
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
