package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main.go/admincontroller"
	"main.go/models"
	"main.go/util"
)

// SignUp page
func GetSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	}
}
func PostSignUp(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		var users models.User

		// Username exist
		if err := db.Where("username=?", username).First(&users).Error; err == nil {
			c.HTML(http.StatusOK, "signup.html", gin.H{"message": "Username already taken"})
			return
		}

		// Email Exist
		if err := db.Where("email=?", email).First(&users).Error; err == nil {
			c.HTML(http.StatusOK, "signup.html", gin.H{"message": "Email already taken"})
			return
		}

		// Hashing the password
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			fmt.Print("Error while hashing the password")
		}

		newUser := models.User{
			Username: username,
			Email:    email,
			Password: string(hash),
		}

		if err := db.Create(&newUser).Error; err != nil {
			fmt.Print("Error while saving to databse")
		}

		c.Redirect(http.StatusSeeOther, "/user/login")
		c.Next()

	}
}

// Login Page

func GetLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func PostLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		// Admin Check
		admin, res := admincontroller.VerifyAdmin(email, password, db)

		if res {
			token, err := util.GenerateJWT(admin.Username)
			if err != nil {
				fmt.Print("Error in setting the token")
			}

			c.SetCookie("token", token, int(24*time.Hour), "/", "localhost", false, true)

			c.Redirect(http.StatusSeeOther, "/admin/userlist")
			return
		}

		var users models.User

		// Email exist
		if err := db.Where("email=?", email).First(&users).Error; err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"message": "Email not exist"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password)); err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"message": "Incorrect Password"})
			return
		}

		token, err := util.GenerateJWT(users.Username)
		if err != nil {
			fmt.Print("Error while generating token")
		}
		c.SetCookie("token", token, int(24*time.Hour), "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/user/home")
		c.Next()
	}
}

func GetHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/user/login")
			return
		}
		claims, err := util.VerifyJWT(token)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/user/login")
			return
		}
		c.HTML(http.StatusOK, "home.html", gin.H{"name": claims.Username})
	}
}

func GetLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
}
