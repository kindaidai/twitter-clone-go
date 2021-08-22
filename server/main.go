package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null;type:varchar(20)"`
	Email    string `gorm:"unique;not null;type:varchar(100)"`
	Password string `gorm:"not null;type:varchar(255)"`
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
		c.HTML(200, "signup.html", gin.H{})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
