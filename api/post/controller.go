package post

import (
	"github.com/gin-gonic/gin"

	"flag"
	"fmt"
	"strconv"
	"challenge3/models"
	"challenge3/database"
	"challenge3/usecase"
)

func GetListPost(postService *usecase.PostService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		flag.Parse()
		_ = flag.Arg(0)
		connection := database.GetDatabase()
	
		p := c.DefaultQuery("page", "1")
		page, err := strconv.Atoi(p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	
		if page <= 0 {
			c.JSON(200, gin.H{
				"message": "Does not exsit page ",
			})
			return
		}
	
		var postList []models.Post
		var offset = (page - 1) *10
		
		connection.Limit(10).Offset(offset).Find(&postList)
		
		c.JSON(200, postList)
	}
	
}

func CreatePost(postService *usecase.PostService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		userID := c.MustGet("userID")
		email := c.MustGet("email")
	
		content := c.PostForm("content")
		var post = models.Post{
			UserID: uint(userID.(float64)),
			Email: email.(string),
			Content: content,
		}
	
		postService.CreatePost(post)
		c.JSON(200, gin.H{
			"message": "Create post successfully",
		})
	}
}

func UpdatePost(postService *usecase.PostService) func(c *gin.Context) {
	return func(c *gin.Context) {

		flag.Parse()
		_ = flag.Arg(0)
	
		role := c.MustGet("role").(string)
		postID1 := c.Param("postID")
		userID := c.MustGet("userID")
		postID, err := strconv.ParseUint(postID1, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	
		postCheck, err := postService.Find(uint(postID))
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Does not exist post",
			})
			return
		}
	
		if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		content := c.PostForm("content")
		postService.UpdatePost(postCheck.ID, content)
	
		c.JSON(200, gin.H{
			"message": "Edit post successfully",
		})
	}
}

func DeletePost(postService *usecase.PostService) func(c *gin.Context) {
	return func(c *gin.Context) {
		
		role := c.MustGet("role").(string)
		postID1 := c.Param("postID")
		userID := c.MustGet("userID")
		postID, err := strconv.ParseUint(postID1, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		postCheck, err := postService.Find(uint(postID))
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Does not exist post",
			})
			return
		}
	
		if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		postService.DeletePost(uint(postID))
		c.JSON(200, gin.H{
			"message": "Delete post successfully",
		})
	}
}