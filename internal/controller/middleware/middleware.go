package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/diana-gemini/drevmass/internal/models"
)

func RequireAuth(c *gin.Context) {

	var tokenString string
	tokenArray := strings.Split(c.GetHeader("Authorization"), " ")

	if len(tokenArray) >= 2 {
		tokenString = tokenArray[1]
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("userID", uint(claims["sub"].(float64)))
		c.Set("roleID", uint(claims["role"].(float64)))
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("userID")

		if userId == 0 {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to find user")
			return
		}

		roleId := c.GetUint("roleID")

		if roleId == 0 {
			models.NewErrorResponse(c, http.StatusBadRequest, "failed to get the role")
			return
		}

		if roleId != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, "forbidden")
			return
		}

		c.Next()
	}
}
