package user

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	
	"bytes"
	"math/rand"
	"time"
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/sha1"
	"challenge3/database"
	"challenge3/models"
	repo "challenge3/repository"
	"challenge3/usecase"
)

var mySigningKey = "pa$$w0rd"
var siteKey = "6Lez19QfAAAAAOD76uihNiKbKv62kL9ap-8KS35_";

func Response(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
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

func GetListUser(userService *usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		userList, err := userService.GetListUser()
	
		if err != nil {
			Response(c, 500, "Database is wrong")
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
	checkPassRepo := repo.NewCheckPassRepo(connection)
	checkPassService := usecase.NewCheckPassService(checkPassRepo)

	var user models.Authen
	clientIP := c.ClientIP()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, err)
		return
	}

	// user.Email = c.PostForm("email")
	// user.Password = c.PostForm("password")
	userAuth, _ := userRepo.Find(user.Email)	
	if userAuth.Email == "" {
		Response(c, 400, "Username or password is incorrect")
		return
	}

	checkPassDoc, errCheck := checkPassService.FindCheck(clientIP, user.Email)
	if errCheck == nil && checkPassDoc.FailedLogin >= 2 {
		captchaKey := user.Captcha
		var body = []byte(`{"secret":"` + siteKey + `", "response":"` + captchaKey + `"}`)
		var r map[string]interface{}

		req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", bytes.NewBuffer(body))
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		if success := r["success"]; success == "false" {
			c.JSON(400, gin.H{})
			return
		}
	}

	if check := CheckPasswordHash(user.Password, userAuth.Password); !check {
		checkPassService.UpdateCheck(clientIP, user.Email)
		Response(c, 400, "Username or password is incorrect")
		return
	} 
	
	tokenString, err := GenerateJWT(&userAuth)
	if err != nil {
		Response(c, 500, "")
		return
	}

	// c.SetCookie("token", tokenString, 150, "/", "localhost", false, true)
	if errCheck == nil {
		checkPassService.DeleteCheck(clientIP, user.Email)
	}
	c.JSON(200, gin.H{
		"token": tokenString,
		"expiredIn": 300,
	})
}

func LogOut(c *gin.Context) {
	c.SetCookie("token", "", 150, "/", "localhost", false, true)
	Response(c, 200, "Successful log out")
}

func Register(userService *usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(500, err)
			return
		}

		email := user.Email //c.PostForm("email")
		name := user.Name //c.PostForm("name")
		password := user.Password //c.PostForm("password")

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	
		if err != nil {
			Response(c, 500, "Cannot generate hash password")
			return
		}

		// userCheck, _ := userService.FindUser(email)
		
		result, err := userService.CreateUser(email, name, string(hashPassword))
		if err != nil {
			Response(c, 400, "Email is already existed")
			return
		}
		c.JSON(200, gin.H{
			"user": result,
		})
	}
}

func CreateUser(userService *usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {	
	
		password := c.PostForm("password")
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	
		if err != nil {
			Response(c, 500, "Cannot generate hash password")
			return
		}
	
		result, err := userService.CreateUser(c.PostForm("email"), c.PostForm("name"), string(hashPassword))
		if err != nil {
			Response(c, 400, "Email is already existed")
			return
		}

		c.JSON(200, gin.H{
			"user": result,
		})
	}
}

func DeleteUser(userService *usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		email := c.Param("userEmail")
	
		err := userService.DeleteUser(email)
		if err != nil {
			Response(c, 400, "Does not exist user")
			return
		}

		c.JSON(204, gin.H{})
	}
}

func UpdateUser(userService *usecase.UserService) func(c *gin.Context) {
	return func(c *gin.Context) {
	
		var user = models.User{
			Email: c.Param("userEmail"),
			Name: c.PostForm("name"),
			Password: c.PostForm("password"),
		}
		
		err := userService.UpdateUser(user)
		if err != nil {
			Response(c, 400, "Do not exist user")
			return
		}
	
		c.JSON(204, gin.H{})
	}

}

func NewRole(roleService *usecase.RoleService) func(c *gin.Context) {
	return func(c *gin.Context) {

		var role models.Role

		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(500, err)
			return
		}
	
		name := role.Name //c.PostForm("name")
		permission := role.Permission //c.PostForm("permission")
	
		// roleCheck, _ := roleService.Find(name)
		
		err := roleService.Create(name, permission)
		if err != nil {
			Response(c, 400, "This role is available")
			return
		}

		Response(c, 200, "Create role successfully")
	}
}

func ChangeRole(userService *usecase.UserService, roleService *usecase.RoleService) func(c *gin.Context) {
	return func(c *gin.Context) {
		connection := database.GetDatabase()
	
		email := c.PostForm("email")
		role := c.PostForm("role")
	
		userCheck, _ := userService.FindUser(email)
		if userCheck.Email == "" {
			Response(c, 400, "Does not exist user")
			return
		}
	
		roleCheck, _ := roleService.Find(role)
		if roleCheck.Name == "" {
			Response(c, 400, "Does not exist role")
			return
		}
	
		userCheck.Role = role
		connection.Save(&userCheck)
		c.JSON(204, gin.H{})
	}

}

func GetMe(userService *usecase.UserService) gin.HandlerFunc {
	return func (c *gin.Context) {
		email := c.MustGet("email").(string)
		user, err := userService.FindUser(email)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, user)
	}
}

func GetUser(userService *usecase.UserService) gin.HandlerFunc {
	return func (c *gin.Context) {
		email := c.Param("userEmail")
		user, err := userService.FindUser(email)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, user)
	}
}

func ChangePass(userService *usecase.UserService) gin.HandlerFunc {
	return func (c *gin.Context) {
		var changePassRequest models.ChangePassRequest
		if err := c.ShouldBindJSON(&changePassRequest); err != nil {
			c.JSON(500, gin.H{})
			return
		}

		user, err := userService.FindUser(changePassRequest.Email)
		if err != nil {
			c.JSON(400, gin.H{})
			return
		}

		if !CheckPasswordHash(changePassRequest.OldPassword, user.Password) {
			c.JSON(400, gin.H{})
			return
		}

		if changePassRequest.NewPassword == changePassRequest.OldPassword {
			c.JSON(400, gin.H{})
			return
		}

		if changePassRequest.ConfirmPassword != changePassRequest.NewPassword {
			c.JSON(400, gin.H{})
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(changePassRequest.NewPassword), 14)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		user.Password = string(hashPassword)
		if err := userService.UpdateUser(user); err != nil {
			c.JSON(500, gin.H{})
			return
		}

		c.JSON(204, gin.H{})
	}

}

func SendResetPassword(userService *usecase.UserService, resetPasswordService *usecase.ResetPasswordService) gin.HandlerFunc {
	return func (c *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(500, gin.H{})
			return
		}
		user, err := userService.FindUser(user.Email)
		if err != nil {
			c.JSON(204, gin.H{})
			return
		}
		randomString := resetPasswordService.GenerateResetPasswordToken(64)
		err = resetPasswordService.StoreToken([]byte(randomString), user.Email)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		c.JSON(200, gin.H{
			"resetPasswordToken": randomString,
		})
	}
}

func ResetPassword(userService *usecase.UserService, resetPasswordService *usecase.ResetPasswordService) gin.HandlerFunc {
	return func (c *gin.Context) {
		token := c.Query("token")
		var (
			changePass models.ChangePassRequest
			user models.PasswordResetToken
			listDoc []models.PasswordResetToken
			matchRow models.PasswordResetToken
		) 
		if err := c.ShouldBindJSON(&changePass); err != nil {
			c.JSON(500, gin.H{})
			return
		}

		user, err := resetPasswordService.FindUser(token)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		listDoc, err = resetPasswordService.FindAllToken(user.Email)
		if err != nil {
			c.JSON(400, gin.H{})
			return
		}

		hasher := sha1.New()
		hasher.Write([]byte(token))
		result := hasher.Sum(nil)
		
		for _, doc := range listDoc {
			if check := bytes.Compare(doc.Token, []byte(result)); check == 0{
				matchRow = doc
				break
			}
		}

		if matchRow.Email == "" && bytes.Compare(matchRow.Token, []byte("")) == 0 {
			c.JSON(400, gin.H{})
			return
		}

		_ = resetPasswordService.DeleteAll(user.Email)
		fmt.Println(matchRow)
		if time.Now().After(matchRow.ExpriedIn) {
			c.JSON(403, gin.H{})
			return
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(changePass.NewPassword), 14)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}

		userUpdate, _ := userService.FindUser(user.Email)
		userUpdate.Password = string(hashPassword)
		if err := userService.UpdateUser(userUpdate); err != nil {
			c.JSON(500, gin.H{})
			return
		}

		c.JSON(204, gin.H{})
	}
}