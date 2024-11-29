package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtkey = []byte("secretkey")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main(){
	r := gin.Default()

	r.POST("/auth", func(c *gin.Context){
		var creds Credentials

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"invalid request"})
			return
		}

		if creds.Username == "admin" && creds.Password == "password" {
			expirationTime := time.Now().Add(1 * time.Hour)
			claims := &jwt.MapClaims{
				"username" : creds.Username,
				"exp" : expirationTime,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtkey)
			if err!=nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
                return
			}
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
		}
	})

	r.Run(":8081")
}