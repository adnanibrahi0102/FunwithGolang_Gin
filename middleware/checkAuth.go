package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/adnanibrahi0102/gin-app/initializers"
	"github.com/adnanibrahi0102/gin-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckAuth(c *gin.Context){ 

	//extract the token from the cookie

	tokenString , err := c.Cookie("jwtToken")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized  request, please login"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse and validate the token

	token , err := jwt.Parse(tokenString , func(token *jwt.Token)(interface{},error){

		if _ , ok := token.Method.(*jwt.SigningMethodHMAC) ; !ok {
			return nil , fmt.Errorf("unexpected signing method : %v" , token.Header["alg"])
		}
		return []byte (os.Getenv("JWT_SECRET")) , nil
	}) 

	if err != nil {
		c.JSON(http.StatusUnauthorized , gin.H{"error": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims , ok  := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized , gin.H{"error": "invalid token"})
		c.Abort()
		return
	}


	// check if the token is expired

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized , gin.H{"error": "token is expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// verify user exists in the database

	user := models.User{}
	initializers.DB.Where("ID = ?" , claims["id"]).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized , gin.H{"error": "user not found"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser" , user)

	c.Next()
}