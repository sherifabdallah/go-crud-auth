package db

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rest-api/models"
)

var DB *gorm.DB

func InitDB() {
    var err error

    // Open the PostgreSQL database
    dsn := "host=localhost user=postgres password=ameame dbname=DB port=5432 sslmode=disable"
    DB, err = gorm.Open(postgres.New(postgres.Config{
        DSN: dsn,
    }), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    err = DB.AutoMigrate(&models.Event{})
    if err != nil {
        panic("failed to migrate database")
    }
}
