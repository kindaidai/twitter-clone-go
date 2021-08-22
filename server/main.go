package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	UserId  int    `gorm:"not null;constrant:OnDelete:CASCADE"`
	User    User
}

type Follow struct {
	gorm.Model
	FollowerId int `gorm:"not null;constrant:OnDelete:CASCADE"`
	Follower   User
	FollowedId int `gorm:"not null;constrant:OnDelete:CASCADE"`
	Followed   User
}

func dbConnect() *gorm.DB {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, dbname)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed DB connect")
	}
	return db
}

func createUser(name string, email string, password string) (*gorm.DB, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 13)
	user := User{Name: name, Email: email, Password: hashedPassword}
	db := dbConnect()
	result := db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

func main() {
	db := dbConnect()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Tweet{}, &Follow{})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "twitter-clone-go",
		})
	})
	router.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{
			"err": "",
		})
	})
	router.POST("/signup", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")
		if _, err := createUser(name, email, password); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err.Error()})
		}
		c.Redirect(302, "/")
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
