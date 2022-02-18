package main

import (
	"os"
	"strconv"
	"strings"

	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	gorm.Model
	Name        string
	Age         int
	Description string
}

func db_init() {
	_, err := os.Stat("test.sqlite3")
	if err == nil {
		os.Remove("test.sqlite3")
	}
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect DB")
	}
	db.AutoMigrate(&Person{})

	// insert seed data before playing
	db.Create(&Person{Name: "admin", Age: 1, Description: os.Getenv("FLAG")})
	db.Create(&Person{Name: "Alice", Age: 20, Description: "Hello, World!"})
	db.Create(&Person{Name: "Bob", Age: 30, Description: "Gorm is really interesting!"})
}

func create(name string, age int, description string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect DB")
	}
	db.Create(&Person{Name: name, Age: age, Description: description})
}

func get_all() []Person {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var people []Person
	// exclude admin record!!
	db.Where("id != 1").Find(&people)
	return people
}

func get_person_by_id(id string) []Person {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	var person []Person
	db.Debug().Find(&person, id)
	return person
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("robots.txt", "./robots.txt")
	db_init()
	r.GET("/", func(c *gin.Context) {
		people := get_all()
		c.HTML(200, "index.tmpl", gin.H{
			"people": people,
		})
	})
	r.POST("/new", func(c *gin.Context) {
		name := c.PostForm("name")
		age, _ := strconv.Atoi(c.PostForm("age"))
		description := c.PostForm("description")
		create(name, age, description)
		c.Redirect(302, "/")
	})
	r.GET("/person", func(c *gin.Context) {
		id := c.Query("id")

		if id == "" {
			c.String(400, "id param is invalid")
			return
		}

		if strings.Contains(id, "1") {
			c.String(401, "you can't read admin record HAHAHA!!")
			return
		}
		person := get_person_by_id(id)

		c.HTML(200, "index.tmpl", gin.H{
			"people": person,
		})

	})
	r.Run(":8888")
}
