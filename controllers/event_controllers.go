package controllers

import (
	"net/http"
	"rest-api/db"
	"rest-api/models"
	"rest-api/services"
	"strconv"
	"github.com/gin-gonic/gin"
)

// Handler function to get all events
func GetEventsController(c *gin.Context) {
	events, err := services.GetEventsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// Handler function to create a new event
func CreateEventController(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := event.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateEventService(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// Handler function to update an existing event by its ID
func UpdateEventController(c *gin.Context) {
	var event models.Event
	id := c.Param("id") // Extract the event ID from the request parameters

	// Convert ID parameter to uint
	eventID, err := strconv.ParseUint(id, 10, 64)

	// Check if id is not valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Find the existing event by its ID
	existingEvent := services.GetEventByIDService(uint(eventID))
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Bind JSON data to the event struct
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the event
	if err := event.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the existing event with the new data
	existingEvent.Name = event.Name
	existingEvent.Description = event.Description
	existingEvent.Location = event.Location
	existingEvent.DateTime = event.DateTime

	// Save the updated event to the database
	if err := services.UpdateEventService(existingEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, existingEvent)
}

// Handler function to delete an event
func DeleteEventController(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := db.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := services.DeleteEventService(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
