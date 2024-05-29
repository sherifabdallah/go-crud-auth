// main.go
package main

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"rest-api/controllers"
	"rest-api/db"
	"rest-api/services"
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

	// ----------------------- Keycloak -----------------------
	// Set up OIDC provider
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/realms/dotnet")
	if err != nil {
		panic(err)
	}

	// Configure OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     "golang",
		ClientSecret: "4y0ShiZu4RiPGtBVOKI7e1GbHSqbnZGM",
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8000/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Token Verifier
	verifier := provider.Verifier(&oidc.Config{ClientID: "golang"})

	// ----------------------- EndKeycloak -----------------------

	// Login route
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL("state"))
	})

	// Secure route
	router.GET("/logedin", services.RequireLogin(verifier), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, you are logged in!"})
	})

	// Callback route
	router.GET("/callback", func(c *gin.Context) {
		ctx := c.Request.Context()
		code := c.Query("code")

		token, err := oauth2Config.Exchange(ctx, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
			return
		}

		idTokenVerifier, err := verifier.Verify(ctx, idToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token: " + err.Error()})
			return
		}

		claims := map[string]interface{}{}
		if err := idTokenVerifier.Claims(&claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save the ID token and claims in the session or any secure storage
		c.JSON(http.StatusOK, gin.H{"message": "Logged in as " + claims["email"].(string)})
	})

	// Start the server
	router.Run(":8000")
}
