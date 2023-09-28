package admincontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main.go/models"
	"main.go/util"
)

func VerifyAdmin(email, password string, db *gorm.DB) (*models.Admin, bool) {

	var admin models.Admin

	if err := db.Where("email=?", email).First(&admin).Error; err != nil {
		return nil, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, false
	}
	return &admin, true

}

func GetUserlist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User

		if err := db.Order("username").Find(&users).Error; err != nil {
			fmt.Print("Error in listing the user")
			return
		}
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
		c.HTML(http.StatusOK, "userlist.html", gin.H{"users": users, "admin": claims.Username})

	}
}

func SearchUser(db *gorm.DB) gin.HandlerFunc {
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

		var users []models.User
		if err := db.Order("username").Find(&users).Error; err != nil {
			fmt.Print("Error in listing the user")
			return
		}
		var searchusers []models.User
		name := c.PostForm("search")
		if name == "" {
			c.Redirect(http.StatusSeeOther, "/admin/userlist")
			return
		}
		username := "%" + name + "%"

		if err := db.Where("username ILIKE ?", username).Find(&searchusers).Error; err != nil {
			c.HTML(http.StatusOK, "userlist.html", gin.H{"admin": claims.Username, "users": users, "searcherror": "There is no user"})
		}

		c.HTML(http.StatusOK, "userlist.html", gin.H{"admin": claims.Username, "users": users, "search": searchusers})
	}
}

func EditUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var users models.User

		if err := db.Where("username=?", username).First(&users).Error; err != nil {
			fmt.Print("Error in finding the user")
			return
		}

		c.HTML(http.StatusOK, "edituser.html", gin.H{"user": users})
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		oldusername := c.Param("username")
		newusername := c.PostForm("username")
		newemail := c.PostForm("email")

		var oldusers models.User

		if err := db.First(&oldusers, "username=?", oldusername).Error; err != nil {
			fmt.Print("Error to find the user")
			return
		}

		var newusers models.User

		// Chechk user already exist

		if err := db.Where("username=?", newusername).Not("username IN (?)", oldusers.Username).First(&newusers).Error; err == nil {
			c.HTML(http.StatusOK, "edituser.html", gin.H{"message": "Username already exist", "user": oldusers})
			return
		}

		// Check email already exist

		if err := db.Where("email=?", newemail).Not("email IN (?)", oldusers.Email).First(&newusers).Error; err == nil {
			c.HTML(http.StatusOK, "edituser.html", gin.H{"message": "Email already exist", "user": oldusers})
			return
		}

		// newpassword := c.PostForm("password")

		oldusers.Username = newusername
		oldusers.Email = newemail

		// if newpassword != "" {
		// 	hash, err := bcrypt.GenerateFromPassword([]byte(newpassword), 12)
		// 	if err != nil {
		// 		fmt.Print("Error in hashing the password")
		// 		return
		// 	}
		// 	users.Password = string(hash)
		// }

		if err := db.Save(&oldusers).Error; err != nil {
			fmt.Print("Error in savig")
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/userlist")

	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		var users models.User

		if err := db.Where("username=?", username).First(&users).Error; err != nil {
			fmt.Print("Error in finding the user")
			return
		}
		if err := db.Delete(&users).Error; err != nil {
			fmt.Print("Error in deleting user")
			return
		}
		c.Redirect(http.StatusSeeOther, "/admin/userlist")
	}
}

func GetCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "createuser.html", nil)
	}
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		//Check the user already exist
		var user models.User

		if err := db.Where("username=?", username).First(&user).Error; err == nil {
			c.HTML(http.StatusOK, "createuser.html", gin.H{"message": "User already exist"})
			return
		}

		// Email check

		if err := db.Where("email=?", email).First(&user).Error; err == nil {
			c.HTML(http.StatusOK, "createuser.html", gin.H{"message": "Email already exist"})
			return
		}

		//Hashing the password

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			fmt.Print("Error in Hasing the password")
		}

		//create new user to the database

		newUser := models.User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := db.Create(&newUser).Error; err != nil {
			fmt.Print("Error in create new user")
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/userlist")
	}
}

func GetCreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "createadmin.html", nil)
	}
}

func CreateAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		var admin models.Admin

		// Check name already exist
		if err := db.Where("username=?", username).First(&admin).Error; err == nil {
			c.HTML(http.StatusOK, "createadmin.html", gin.H{"message": "Name is already exist"})
			return
		}

		// Check email already exist
		if err := db.Where("email=?", email).First(&admin).Error; err == nil {
			c.HTML(http.StatusOK, "createadmin.html", gin.H{"message": "Name is already exist"})
			return
		}

		// Hashing the password
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			fmt.Print("Error in hashing the password")
			return
		}

		// Adding to databse

		newUser := models.Admin{
			Username: username,
			Email:    email,
			Password: string(hashPassword),
		}

		if err := db.Create(&newUser).Error; err != nil {
			fmt.Print("Error in hashing the password")
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/userlist")
	}
}
