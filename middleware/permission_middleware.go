package middleware

import (
	"challenge3/database"
	repo "challenge3/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

func NeedPermission(permit string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role")
		connection := database.GetDatabase()
		roleRepo := repo.NewRoleRepo(connection)

		permission, _ := roleRepo.Find(role.(string))

		if ok := strings.Contains(permission.Permission, permit); !ok {

			c.AbortWithStatusJSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}

		c.Next()
	}
}

func NeedRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCheck := c.MustGet("role").(string)
		if role != roleCheck {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "Not Authorized",
			})
			return
		}
		
		c.Next()
	}
}