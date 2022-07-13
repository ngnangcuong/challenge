package post

import (
	"github.com/gin-gonic/gin"

	"flag"
	"fmt"
	"time"
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
			c.JSON(400, gin.H{
				"message": "Does not exsit page ",
			})
			return
		}
	
		var postList []models.Post
		var offset = (page - 1) *10
		
		connection.Order("create_at desc").Limit(10).Offset(offset).Find(&postList)
		
		c.JSON(200, postList)
	}
	
}

func CreatePost(postService *usecase.PostService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		userID := c.MustGet("userID")
		email := c.MustGet("email")
	
		// content := c.PostForm("content")
		var post models.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(400, err.Error())
			return
		}
		// var post = models.Post{
		// 	UserID: uint(userID.(float64)),
		// 	Email: email.(string),
		// 	Content: content.content,
		// 	Create_At: time.Now(),
		// }
		post.UserID = uint(userID.(float64))
		post.Email = email.(string)
		post.Create_At = time.Now()
	
		newPost, err := postService.CreatePost(post)
		if err != nil {
			c.JSON(500, err)
		}

		c.JSON(201, newPost)
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
			c.JSON(400, gin.H{
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
	
		// content := c.PostForm("content")
		var newPost models.Post
		err = c.ShouldBindJSON(&newPost)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		postService.UpdatePost(postCheck.ID, newPost.Content)
	
		c.JSON(204, gin.H{})
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
			c.JSON(400, gin.H{
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
		c.JSON(204, gin.H{})
	}
}

func SearchPost(postService *usecase.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Param("keyword")
		postList, err := postService.SearchPosts(keyword)
		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, postList)
	}
}

func GetPost(postService *usecase.PostService) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("postId")
		postId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		post, err := postService.Find(uint(postId))
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, post)
	}
}

func FindPostByEmail(postService *usecase.PostService) gin.HandlerFunc {
	return func (c *gin.Context) {
		userEmail := c.Param("userEmail")
		postList, err := postService.FindByEmail(userEmail)

		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, postList)
	}
}
