package main

import (
	"SamkoOfMraz/config"
	"SamkoOfMraz/db"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	config.SetupRouter(router)

	fmt.Println("Server is running")
	if db.IsDatabaseRunning() {
		fmt.Println("Database is running")
	} else {
		fmt.Println("Database is not running")
	}
	err := router.Run("localhost:9090")

	if err != nil {
		fmt.Println("server is not working try restarting server or contact me on emanuel.huzvar@kosickaakademia.sk")
		fmt.Println(err)
	}

}
