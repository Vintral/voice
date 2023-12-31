package main

import (
	"api/models"
	"api/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/otel/trace"
)

var key = []byte("aBcDeFgHiJkLmNOpQrStUvWxYz")
var tracer trace.Tracer

func Authorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if (len(authorization) == 0) || (strings.Index(authorization, " ") == -1) {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
		} else {
			tokenString := strings.Split(authorization, " ")[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return key, nil
			})

			if err != nil {
				fmt.Println(err)

				c.String(http.StatusUnauthorized, "Unauthorized")
				c.Abort()
			} else {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					c.Set("email", claims["username"])
					c.Next()
				} else {
					fmt.Println(err)

					c.String(http.StatusUnauthorized, "Unauthorized")
					c.Abort()
				}
			}
		}
	}
}

func AdminAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := utils.GetEmailFromContext(c)

		db, err := models.Database(false)
		if err != nil {
			panic(err)
		}

		var user models.User
		db.Where("email = ?", email).First(&user)

		if user.Admin {
			c.Next()
		} else {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
		}
	}
}
