package main

import (
	"log"
	"net/http"

	"github.com/adnanibrahi0102/gin-app/controllers"
	"github.com/adnanibrahi0102/gin-app/initializers"
	"github.com/adnanibrahi0102/gin-app/middleware"
	"github.com/adnanibrahi0102/gin-app/models"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {
    err := initializers.DB.AutoMigrate(&models.User{}, &models.Post{})
	
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Migrations completed successfully")

	router := gin.Default()
	router.GET("/ping", ping)
	router.POST("/auth/register" , controllers.RegisterUser)
	router.POST("/auth/login" , controllers.LoginUser)
	router.GET("/user/me" , middleware.CheckAuth , controllers.GetUser)
	router.POST("/auth/logout" , controllers.LogoutUser)

	//post routes
	router.POST("/posts/create" ,middleware.CheckAuth , controllers.CreatePost)
	router.DELETE("/posts/delete/:id" , middleware.CheckAuth , controllers.DeletePostByID)
	router.GET("/posts/allposts" , middleware.CheckAuth , controllers.GetAllPostsByUserID)
	router.Run(":9090") 
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from golang server ",
	})
}