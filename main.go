package main

import (
	"fmt"
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a simple user structure
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Task struct {
	title string
	desc  string
}

type Users []User

var users = Users{}

func (*Users) newUser(id string, username string, email string) (*Users, error) {
	newUser := &User{ID: id, Username: username, Email: email}
	users = append(users, *newUser)
	return &users, nil
}


func main() {
	uri := "mongodb+srv://ADMIN:pass1234@db.brs9ce0.mongodb.net/?retryWrites=true&w=majority"

	ctx := context.Background()
	client,err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	// Create a client with appropriate options
	if err != nil {
	    fmt.Println("Error connecting to MongoDB:", err)
	    return
	}

	// Ping the server to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
	    fmt.Println("Error pinging MongoDB:", err)
	    return
	}

	// Use the client for database operations
	fmt.Println("Successfully connected to MongoDB!")

	// Remember to close the client when done
	defer func() {
	    if err := client.Disconnect(ctx); err != nil {
	        fmt.Println("Error disconnecting from MongoDB:", err)
	    }
	}()
	router := gin.Default()

	// Define routes for the RESTful API
	router.GET("/users", GetUsers)
	router.GET("/users/:id", GetUser)
	router.POST("/users", CreateUser)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	router.Run(fmt.Sprintf(":%d", port))

}

// GetUsers handles GET requests to retrieve all users
func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

// GetUser handles GET requests to retrieve a specific user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// CreateUser handles POST requests to create a new user
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users = append(users, user)
	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT requests to update an existing user by ID
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			var updatedUser User
			if err := c.ShouldBindJSON(&updatedUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			users[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

// DeleteUser handles DELETE requests to delete a user by ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}
