package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/adnanibrahi0102/gin-app/initializers"
	"github.com/adnanibrahi0102/gin-app/models"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	postInput := models.PostInput{}
	if err := c.ShouldBind(&postInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join("uploads", "images", fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save image"})
		return
	}

	// Retrieve user id from context
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to retrieve user from context"})
		return
	}

	// Create the post
	newPost := models.Post{
		Title:     postInput.Title,
		Content:   postInput.Content,
		Image:     savePath,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := initializers.DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	// Fetch the post with the User relationship preloaded
	var createdPost models.Post
	if err := initializers.DB.Preload("User").First(&createdPost, newPost.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post created successfully",
		"post":    createdPost,
	})
}

// delete post

func DeletePostByID(c *gin.Context) {
	// get post id from url
	postID := c.Param("id")

	// check if post exists
	var post models.Post
	if err := initializers.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// check if the user is the owner of the post
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to retrieve user from context"})
		return
	}

	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to delete this post"})
		return
	}
	// delete post
	if err := initializers.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})

}

// get all posts of a user

func GetAllPostsByUserID(c *gin.Context) {
	//ALGORITHM
	// 1. Get current user from context
	currentUser, exists := c.Get("currentUser")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := currentUser.(models.User)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to retrieve user from context"})
		return
	}
	// 2. Get all posts of the user
	var posts []models.Post
	if err := initializers.DB.Preload("User").Where("user_id = ?", user.ID).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"allPosts": posts, "message": "posts fetched successfully", "totalPosts": len(posts)})

}
