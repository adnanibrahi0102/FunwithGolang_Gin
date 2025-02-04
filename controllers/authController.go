package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adnanibrahi0102/gin-app/initializers"
	"github.com/adnanibrahi0102/gin-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	authInput := models.AuthInput{}

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser := models.User{}

	initializers.DB.Where("username = ? AND email = ?", authInput.Username, authInput.Email).First(&existingUser)

	if existingUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already in use"})
		return

	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("hashPassword", hashPassword)

	user := models.User{
		Username: authInput.Username,
		Email:    authInput.Email,
		Password: string(hashPassword),
	}

	initializers.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})

}
func LoginUser(c *gin.Context) {
	authInput := models.AuthInput{}

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser := models.User{}

	initializers.DB.Where("username = ? AND email = ?", authInput.Username, authInput.Email).Find(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       existingUser.ID,
		"email":    existingUser.Email,
		"username": existingUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}
	c.SetCookie("jwtToken", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User Logged In Successfully",
	})
}

func GetUser(c *gin.Context) {
	user, _ := c.Get("currentUser")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func LogoutUser(c *gin.Context) {
	c.SetCookie("jwtToken", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User Logged Out Successfully",
	})
}
