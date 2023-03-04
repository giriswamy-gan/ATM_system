package main

import (
	"atm-system/configs"
	"atm-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(r)

	r.Run("localhost:9003")
}
