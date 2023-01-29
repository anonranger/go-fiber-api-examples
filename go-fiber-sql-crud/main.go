package main

import (
	db "go-fiber-crud/database"
	router "go-fiber-crud/routes"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	// Using Fiber as the web framework to start the server and handle the routes
	server := fiber.New()

	// Initialize the database connection
	db.IniDatabase()

	// Add the routes to the server
	router.AddRoutes(server)

	// Start the server on port 3000
	server.Listen(":3000")

	// Close the database connection when the server is stopped
	defer db.DBConn.Close()
}

//pasting router code in main file for easier understanding

package routes

import (
	handlers "go-fiber-crud/pkg/users/handlers"

	"github.com/gofiber/fiber/v2"
)

// Set the routes for the server
const (
	GetUsersRoute   = "/api/v1/users"
	GetUserRoute    = "/api/v1/user/:id"
	CreateUserRoute = "/api/v1/user"
	UpdateUserRoute = "/api/v1/user/:id"
	DeleteUserRoute = "/api/v1/user/:id"
)

// AddRoutes adds the routes to the server
func AddRoutes(s *fiber.App) {
	s.Get(GetUsersRoute, handlers.GetUsers)
	s.Get(GetUserRoute, handlers.GetUser)
	s.Post(CreateUserRoute, handlers.CreateUser)
	s.Put(UpdateUserRoute, handlers.UpdateUser)
	s.Delete(DeleteUserRoute, handlers.DeleteUser)
}

//database and models code 

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBConn is the database connection object
var (
	DBConn *gorm.DB
)

// IniDatabase initializes the database connection
func IniDatabase() {

	var err error

	// Open a database connection and save the reference to `DBConn` object
	DBConn, err = gorm.Open("sqlite3", "./db/test.db")
	if err != nil {
		panic("failed to connect database")
	}
	return
	
}

//models code

package models

// User is the model for the users table in the database
type User struct {
	ID    int64  `json:"id" gorm:"primary_key;auto_increment"`
	Fname string `json:"fname" gorm:"type:varchar(250)"`
	Lname string `json:"lname" gorm:"type:varchar(250)"`
	Email string `json:"email" gorm:"type:varchar(250)"`
	Phone int64  `json:"phone" gorm:"type:int"`
}

// TableName sets the insert table name for this struct type
func (b *User) TableName() string {
	return "users"
}


// pasting handler code in same file for easy understanding

package users

import (
	"go-fiber-crud/database"
	models "go-fiber-crud/pkg/users/models"

	"github.com/gofiber/fiber/v2"
)

// GetUsers returns all the users in the database
func GetUsers(c *fiber.Ctx) error {

	// Get the database connection object
	db := database.DBConn
	// Create an array of users
	var users []models.User
	// Get all the users from the database and store it in the array
	db.Find(&users)
	// Return the array of users as JSON
	return c.JSON(users)

}

// GetUser returns a single user from the database
func GetUser(c *fiber.Ctx) error {

	// Get the id from the request parameters
	id := c.Params("id")
	// Get the database connection object
	db := database.DBConn
	// Create a user object
	var user models.User
	// Get the user from the database and store it in the user object
	db.Find(&user, id)
	// Return the user object as JSON
	return c.JSON(user)
}

// CreateUser creates a new user in the database
func CreateUser(c *fiber.Ctx) error {

	// Get the database connection object
	db := database.DBConn
	// Create a user object
	user := new(models.User)
	// Parse the request body and store it in the user object
	if err := c.BodyParser(user); err != nil {
		return err
	}
	// Save the user object in the database
	db.Create(&user)
	return c.JSON("User created successfully!")
}

// UpdateUser updates a user in the database
func UpdateUser(c *fiber.Ctx) error {

	// Get the id from the request parameters
	id := c.Params("id")
	// Get the database connection object
	db := database.DBConn
	// Create a user object
	var user models.User
	// Get the user from the database and store it in the user object
	db.First(&user, id)
	// Check if the user exists
	if user.Fname == "" {
		return c.Status(500).SendString("User not found with ID")
	}
	// Parse the request body and store it in the user object
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	// Save the user object in the database
	db.Save(&user)

	return c.JSON("User updated successfully!")
}

func DeleteUser(c *fiber.Ctx) error {

	// Get the id from the request parameters
	id := c.Params("id")
	// Get the database connection object
	db := database.DBConn
	// Create a user object
	var user models.User
	// Get the user from the database and store it in the user object
	db.First(&user, id)
	// Check if the user exists
	if user.Fname == "" {
		return c.Status(500).SendString("User not found with ID")
	}
	// Delete the user from the database
	db.Delete(&user)

	return c.JSON("User deleted successfully!")
}

