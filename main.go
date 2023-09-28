package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/admincontroller"
	"main.go/controller"
	"main.go/middleware"
	"main.go/models"
)

func main() {

	// Router
	router := gin.Default()

	// Loading HTML
	router.LoadHTMLGlob("template/*.html")

	// Connecting to Database
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env")
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print("Error to connect the database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Admin{})

	// User Part
	user := router.Group("/user")
	user.Use(middleware.DisableCaching())

	// Login, SignUp, Home page setup
	user.GET("/signup", middleware.LoginAuthentication(db), controller.GetSignUp())

	user.POST("/signup", controller.PostSignUp(db))

	user.GET("/login", middleware.LoginAuthentication(db), controller.GetLogin())

	user.POST("/login", controller.PostLogin(db))

	user.GET("/home", middleware.HomeAuthentication(), controller.GetHome())

	user.GET("/logout", controller.GetLogout())

	// Admin Part
	admin := router.Group("/admin")
	admin.Use(middleware.AdminAuthentication(db), middleware.DisableCaching())

	admin.GET("/userlist", admincontroller.GetUserlist(db)) //Show the admin page

	admin.POST("/userlist", admincontroller.SearchUser(db)) //Search the users

	admin.GET("/:username/edituser", admincontroller.EditUser(db)) //Show the edit page

	admin.POST("/:username/edituser", admincontroller.UpdateUser(db)) //Update user

	admin.POST("/:username/deleteuser", admincontroller.DeleteUser(db)) //Delete user

	admin.GET("/createuser", admincontroller.GetCreateUser()) //Show create page

	admin.POST("/createuser", admincontroller.CreateUser(db)) //Create user

	admin.GET("/createadmin", admincontroller.GetCreateAdmin())

	admin.POST("/createadmin", admincontroller.CreateAdmin(db))

	// Logout part

	// Router Run
	router.Run()
}
