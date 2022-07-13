package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	"challenge3/api/user"
	"challenge3/api/post"
	"challenge3/database"
	"challenge3/middleware"
	repo "challenge3/repository"
	es "challenge3/repository/elasticsearch"
	"challenge3/usecase"
)

func NewOpenAPIMiddleware() gin.HandlerFunc {
	validator := middleware.OpenapiInputValidator("./openapi.yaml")
	return validator
}

func InitRoute(router *gin.Engine) {
	validator := NewOpenAPIMiddleware()

	connection := database.GetDatabase()
	esClient := database.GetESClient()

	userRepo := repo.NewUserRepo(connection)
	userService := usecase.NewUserService(userRepo)

	postRepo := repo.NewPostRepo(connection)
	postSearchRepo := es.NewPostSearchRepo(esClient)
	postService := usecase.NewPostService(postRepo, postSearchRepo)

	roleRepo := repo.NewRoleRepo(connection)
	roleService := usecase.NewRoleService(roleRepo)

	resetPasswordRepo := repo.NewResetPasswordRepo(connection)
	resetPasswordService := usecase.NewResetPasswordService(resetPasswordRepo)

	router.Use(middleware.SetupCors()) 

	router.POST("/user/login", user.LogIn)
	router.GET("/user/logout", user.LogOut)
	router.POST("/user/register", user.Register(userService))
	router.PUT("/user/changePass", user.ChangePass(userService))
	router.POST("/user/resetPassword", user.SendResetPassword(userService, resetPasswordService))
	router.PUT("/user/resetPassword", user.ResetPassword(userService, resetPasswordService))
	router.GET("/post", post.GetListPost(postService))
	router.GET("/post/search/:keyword", post.SearchPost(postService))
	
	userRoute := router.Group("/user")
	{
		userRoute.Use(validator)
		userRoute.Use(middleware.SetupCors()) 
		userRoute.Use(middleware.Authorized())
		
		userRoute.POST("/create-user", middleware.NeedPermission("c"), user.CreateUser(userService))
		userRoute.DELETE("/delete-user/:userEmail", middleware.NeedPermission("d"), user.DeleteUser(userService))
		userRoute.PATCH("/update-user/:userEmail", middleware.NeedPermission("u"), user.UpdateUser(userService))
		userRoute.PUT("/change-role", middleware.NeedRole("admin"), user.ChangeRole(userService, roleService))
		userRoute.POST("/new-role", middleware.NeedRole("admin"), user.NewRole(roleService))
		userRoute.GET("/me", user.GetMe(userService))
		userRoute.GET("/:userEmail", user.GetUser(userService))
		userRoute.GET("/", user.GetListUser(userService))
	}

	postRoute := router.Group("/post")
	{
		postRoute.Use(validator)
		postRoute.Use(middleware.SetupCors()) 
		postRoute.Use(middleware.Authorized())
		
		postRoute.POST("/create", post.CreatePost(postService))
		postRoute.DELETE("/delete/:postID", post.DeletePost(postService))
		postRoute.PUT("/update/:postID", post.UpdatePost(postService))
		postRoute.GET("/user/:userEmail", post.FindPostByEmail(postService))
		postRoute.GET("/:postId", post.GetPost(postService))
		// postRoute.GET("/search/:keyword", post.SearchPost(postService))
		// postRoute.GET("/", post.GetListPost(postService))
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

