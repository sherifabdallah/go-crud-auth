// main.go
package main

import (
	"rest-api/controllers"
	"rest-api/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.InitDB()

	// Create a new Gin router
	router := gin.Default()

	// Define routes
	router.GET("/events", controllers.GetEventsController)
	router.POST("/events", controllers.CreateEventController)
	router.PUT("/events/:id", controllers.UpdateEventController)
	router.DELETE("/events/:id", controllers.DeleteEventController)

	// Start the server
	router.Run(":8000")
}
