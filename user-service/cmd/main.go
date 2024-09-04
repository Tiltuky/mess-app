package main

import (
	"github.com/joho/godotenv"
	"log"
	"user-service/app"
)

// @title dating-app user-service API
// @version         1.0
// @description     A users management service API in Go using Gin framework.
// @securityDefinitions.basic  BasicAuth
// @host      localhost:8080
// @BasePath  /api/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server := app.NewAppObj()
	server.AddRoutes()
	server.RunApp()

}
