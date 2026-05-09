package middleware

import (
	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/response"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Unauthenticated")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Unauthenticated")
			c.Abort()
			return
		}

		token, err := jwtService.ValidateToken(parts[1])
		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Unauthenticated")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Unauthenticated")
			c.Abort()
			return
		}

		userIDStr := claims["user_id"].(string)
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthenticated")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
