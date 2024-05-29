package services

import (
	"rest-api/db"
	"rest-api/models"
)

// GetEvents retrieves all events from the database
func GetEventsService() ([]models.Event, error) {
    var events []models.Event
    result := db.DB.Find(&events)
    if result.Error != nil {
        return nil, result.Error
    }
    return events, nil
}

// CreateEvent creates a new event in the database
func CreateEventService(event *models.Event) error {
    result := db.DB.Create(event)
    return result.Error
}

// UpdateEvent updates an existing event in the database
func UpdateEventService(event *models.Event) error {
    result := db.DB.Save(event)
    return result.Error
}

// DeleteEvent deletes an event from the database
func DeleteEventService(event *models.Event) error {
    result := db.DB.Delete(event)
    return result.Error
}

// GetEventByID retrieves an event from the database by its ID
func GetEventByIDService(id uint) *models.Event {
    var event models.Event
    result := db.DB.First(&event, id)
    if result.Error != nil {
        // Handle the case where the event is not found
        return nil
    }
    return &event
}