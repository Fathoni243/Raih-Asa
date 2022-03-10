package main

import (
	"fmt"
	"raih-asa/beasiswa"
	"raih-asa/config"
	"raih-asa/lomba"
	"raih-asa/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB
var r *gin.Engine

func InitGin() {
	r = gin.Default()
}

func StartServer() error {
	return r.Run()
}

func main() {

	db, err := config.InitDB()
	if err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}

	InitGin()

	user.InitRouter(db, r)
	beasiswa.InitRouter(db, r)
	lomba.InitRouter(db, r)

	router := gin.Default()
	router.Use(cors.Default())

	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
