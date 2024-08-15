package main

import (
	"github.com/gin-gonic/gin"
	"rest-api/controllers"
	"rest-api/db"
)

func main() {

	// Initialize the database
	db.InitDB()

	// Create a new Gin router
	router := gin.Default()

	// Authentication route
	router.POST("/login", controllers.LoginHandler)

	// Define routes with authentication middleware
	protected := router.Group("/")
	protected.Use(controllers.AuthMiddleware())
	{
		protected.GET("/events", controllers.GetEventsController)
		protected.POST("/events", controllers.CreateEventController)
		protected.PUT("/events/:id", controllers.UpdateEventController)
		protected.DELETE("/events/:id", controllers.DeleteEventController)
	}

	// Start the server
	router.Run(":8000")
}
