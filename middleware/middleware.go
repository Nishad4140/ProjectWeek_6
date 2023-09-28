package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/models"
	"main.go/util"
)

func AdminAuthentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		auth := c.GetHeader("Authorization")
		if auth != "" {
			token = auth
		} else {
			cookie, err := c.Request.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			c.Redirect(http.StatusSeeOther, "/user/login")
			c.Abort()
			return
		}

		claims, err := util.VerifyJWT(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/user/login")
			c.Abort()
			return
		}

		username := claims.Username
		var user models.Admin

		if err := db.Where("username=?", username).First(&user).Error; err != nil {
			c.Redirect(http.StatusSeeOther, "/user/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoginAuthentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err == nil {
			claims, err := util.VerifyJWT(token)
			if err != nil {
				fmt.Print("Error while verifying jwt")
			}
			var admin models.Admin
			if err := db.Where("username=?", claims.Username).First(&admin).Error; err == nil {
				c.Redirect(http.StatusSeeOther, "/admin/userlist")
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusSeeOther, "/user/home")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func HomeAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		auth := c.GetHeader("Authorization")
		if auth != "" {
			token = auth
		} else {
			cookie, err := c.Request.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			c.Redirect(http.StatusSeeOther, "/user/login")
			c.Abort()
			return
		}

		_, err := util.VerifyJWT(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/user/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func DisableCaching() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}
