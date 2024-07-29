package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AhishekOza/E-commerce/database"
	"github.com/AhishekOza/E-commerce/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Instialize the fiber
	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading environment file")
	}

	PORT := os.Getenv("PORT")

	// IF .env file fails to load then we need to set a default port for the server to listen on .
	if PORT == "" {
		PORT = "8000"
	}

	database.ConnectDB()

	routes.Setup(app)

	// start server
	fmt.Printf("Sever Started listening on %s", PORT)
	app.Listen(":" + PORT)

}
