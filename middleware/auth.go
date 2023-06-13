package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Account struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

var JWT_KEY = []byte("dbceria")

type JWT struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var DB *gorm.DB


func MiddlewareAdmin(c *gin.Context) {
	tokenString, err := c.Cookie("gin_cookie")
	if err != nil {
		c.AbortWithStatus(401)
		c.JSON(401, gin.H{"message": "Silahkan Login"})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return JWT_KEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(401)
			c.JSON(401, gin.H{"message": "Silahkan Login Kembali"})
		} 
		var account Account
		DB.First(&account, claims["Id"])
		if account.Id == 0 {
			c.AbortWithStatus(401)
		}
		c.Set("account", account)
		c.Next()
	} else {
		c.AbortWithStatus(401)
	}

}