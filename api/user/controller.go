package user

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	
	"time"
	"challenge3/database"
	"challenge3/models"
	repo "challenge3/repository"
	"challenge3/usecase"
)

var mySigningKey = "pa$$w0rd"

func Response(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func CheckPasswordHash(password string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func GenerateJWT(userAuth *models.User) (string, error) {
	var secretkey = []byte(mySigningKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = userAuth.Email
	claims["userID"] = userAuth.ID
	claims["role"] = userAuth.Role
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(secretkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetListUser(userService usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		if check := c.MustGet("isLogin").(bool); !check {
			c.JSON(200, gin.H{
				"message": "Not Log in yet",
			})
			return
		}
	
		userList, err := userService.GetListUser()
	
		if err != nil {
			Response(c, 200, "Database is wrong")
			return
		}
	
	
		c.HTML(200, "listUser.tmpl", gin.H{
			"userList": userList,
		})
	}
}

func LogIn(c *gin.Context) {
	connection := database.GetDatabase()
	userRepo := repo.NewUserRepo(connection)

	email := c.PostForm("email")
	password := c.PostForm("password")

	userAuth, _ := userRepo.Find(email)
	if userAuth.Email == "" {
		Response(c, 200, "Not User")
		return
	}

	if check := CheckPasswordHash(password, userAuth.Password); !check {
		Response(c, 200, "Password is not correct")
		return
	} 
	
	tokenString, err := GenerateJWT(&userAuth)
	if err != nil {
		Response(c, 200, "Cannot generate jwt token")
		return
	}

	// c.SetCookie("token", tokenString, 150, "/", "localhost", false, true)
	Response(c, 200, tokenString)
}

func LogOut(c *gin.Context) {
	c.SetCookie("token", "", 150, "/", "localhost", false, true)
	Response(c, 200, "Successful log out")
}

func Register(userService usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		email := c.PostForm("email")
		name := c.PostForm("name")
		password := c.PostForm("password")
	
		userCheck, _ := userService.FindUser(email)
		if userCheck.Email != "" {
			Response(c, 200, "Email is already existed")
			return
		}
	
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	
		if err != nil {
			Response(c, 200, "Cannot generate hash password")
			return
		}
	
		userService.CreateUser(email, name, string(hashPassword))
		Response(c, 200, "Create user successfully")
	}
}

func CreateUser(userService usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {	
		if check := c.MustGet("isLogin").(bool); !check {
			c.JSON(200, gin.H{
				"message": "Not Log in yet",
			})
			return
		}
	
		if permit := c.MustGet("Permission").(bool); !permit {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		password := c.PostForm("password")
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	
		if err != nil {
			Response(c, 200, "Cannot generate hash password")
			return
		}
	
		userService.CreateUser(c.PostForm("email"), c.PostForm("name"), string(hashPassword))
	
		Response(c, 200, "Create user successfully")
	}
}

func DeleteUser(userService usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		if check := c.MustGet("isLogin").(bool); !check {
			c.JSON(200, gin.H{
				"message": "Not Log in yet",
			})
			return
		}
	
		if permit := c.MustGet("Permission").(bool); !permit {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		email := c.Param("userEmail")
	
		err := userService.DeleteUser(email)
		if err != nil {
			Response(c, 200, "Does not exist user")
			return
		}

		Response(c, 200, "Delete user successfully")
	}
}

func UpdateUser(userService usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		if check := c.MustGet("isLogin").(bool); !check {
			c.JSON(200, gin.H{
				"message": "Not Log in yet",
			})
			return
		}
	
		if permit := c.MustGet("Permission").(bool); !permit {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		var user = models.User{
			Email: c.Param("userEmail"),
			Name: c.PostForm("name"),
			Password: c.PostForm("password"),
		}
		
		err := userService.UpdateUser(user)
		if err != nil {
			Response(c, 200, "Do not exist user")
			return
		}
	
		Response(c, 200, "Update user successfully")
	}

}

func NewRole(roleService usecase.RoleService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		if check := c.MustGet("isLogin").(bool); !check {
			Response(c, 200, "Not Log in yet")
			return
		}
	
		if permit := c.MustGet("Permission").(bool); !permit {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		name := c.PostForm("name")
		permission := c.PostForm("permission")
	
		roleCheck, _ := roleService.Find(name)
		if roleCheck.Name != "" {
			Response(c, 200, "This role is available")
			return
		}
	
		roleService.Create(name, permission)
		Response(c, 200, "Create role successfully")
	}
}

func ChangeRole(userService usecase.UserService, roleService usecase.RoleService) func(c *gin.Context) {
	return func(c *gin.Context) {
		connection := database.GetDatabase()
		if check := c.MustGet("isLogin").(bool); !check {
			Response(c, 200, "Not Log in yet")
			return
		}
	
		if permit := c.MustGet("Permission").(bool); !permit {
			c.JSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
	
		email := c.PostForm("email")
		role := c.PostForm("role")
	
		userCheck, _ := userService.FindUser(email)
		if userCheck.Email == "" {
			Response(c, 200, "Does not exist user")
			return
		}
	
		roleCheck, _ := roleService.Find(role)
		if roleCheck.Name == "" {
			Response(c, 200, "Does not exist role")
			return
		}
	
		userCheck.Role = role
		connection.Save(&userCheck)
		Response(c, 200, "Change role successfully")
	}

}