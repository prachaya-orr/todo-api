package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func AccessToken(signature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
			Issuer:    "prachaya-orr",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		})

		ss, err := token.SignedString([]byte(signature))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": ss,
		})
	}
}
