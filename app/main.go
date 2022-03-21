package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	"challenge3/api/user"
	"challenge3/api/post"
	"challenge3/database"
	"challenge3/middleware"
	repo "challenge3/repository"
	"challenge3/usecase"
)

func NewOpenAPIMiddleware() gin.HandlerFunc {
	validator := middleware.OpenapiInputValidator("./openapi.yaml")
	return validator
}

func InitRoute(router *gin.Engine) {
	validator := NewOpenAPIMiddleware()

	connection := database.GetDatabase()
	userRepo := repo.NewUserRepo(connection)
	userService := usecase.NewUserService(userRepo)

	postRepo := repo.NewPostRepo(connection)
	postService := usecase.NewPostService(postRepo)

	roleRepo := repo.NewRoleRepo(connection)
	roleService := usecase.NewRoleService(roleRepo)
	
	userRoute := router.Group("/user")
	{
		userRoute.Use(validator)
		userRoute.Use(middleware.Authorized())

		userRoute.POST("/login", user.LogIn)
		userRoute.GET("/logout", user.LogOut)
		userRoute.POST("/register", user.Register(userService))
		userRoute.POST("/create-user", middleware.NeedPermission("c"), user.CreateUser(userService))
		userRoute.DELETE("/delete-user/:userEmail", middleware.NeedPermission("d"), user.DeleteUser(userService))
		userRoute.PATCH("/update-user/:userEmail", middleware.NeedPermission("u"), user.UpdateUser(userService))
		userRoute.PUT("/change-role", middleware.NeedRole("admin"), user.ChangeRole(userService, roleService))
		userRoute.POST("/new-role", middleware.NeedRole("admin"), user.NewRole(roleService))
		userRoute.GET("/", user.GetListUser(userService))
	}

	postRoute := router.Group("/post")
	{
		postRoute.Use(validator)
		postRoute.Use(middleware.Authorized())

		postRoute.POST("/create", post.CreatePost(postService))
		postRoute.DELETE("/delete/:postID", post.DeletePost(postService))
		postRoute.PUT("/update/:postID", post.UpdatePost(postService))
		postRoute.GET("/", post.GetListPost(postService))
	}
}

func InitAPI() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	InitRoute(router)
	router.Run(":3000")
}

func main() {
	database.InitMigration()
	InitAPI()
}

