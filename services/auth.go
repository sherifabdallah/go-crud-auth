package services

import (
	"github.com/coreos/go-oidc"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Require Login
func RequireLogin(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
            c.Abort()
            return
        }

        idToken := authHeader[len("Bearer "):]
        ctx := c.Request.Context()
        token, err := verifier.Verify(ctx, idToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify ID token: " + err.Error()})
            c.Abort()
            return
        }

        // Optionally, you can check additional claims in the token
        claims := map[string]interface{}{}
        if err := token.Claims(&claims); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        // Pass claims to the context for further use in handlers
        c.Set("claims", claims)

        // Call the next handler
        c.Next()
    }
}
