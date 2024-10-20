package main

import (
	"encoding/json"
	"fmt"
	"go_userlist/repository"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the repository
	userRepo, err := repository.NewPostgresUserRepository(nil)
	if err != nil {
		fmt.Println("Error initializing repository:", err)
	}
	defer userRepo.Close()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Define routes
	r.GET("/users", func(c *gin.Context) { getAllUsersHandler(c, *userRepo) }) // Pass userRepo which implements UserRepository
	r.GET("/users/:id", func(c *gin.Context) {
		// Extract and convert the userId from the URL
		idParam := c.Param("id")
		userId, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Fetch the user by ID
		user, err := userRepo.GetUserByID(userId)
		if err != nil {
			if err == repository.ErrUserNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Return the user details as JSON
		c.JSON(http.StatusOK, user)
	})
	r.POST("/users", func(c *gin.Context) { createUserHandler(c, *userRepo) })     // Pass userRepo
	r.PUT("/users/:id", func(c *gin.Context) { updateUserHandler(c, *userRepo) })  // Pass userRepo
	r.PATCH("/users/:id", func(c *gin.Context) { patchUserHandler(c, *userRepo) }) // Pass userRepo
	r.DELETE("/users/:id", func(c *gin.Context) { deleteUserHandler(c, *userRepo) })
	// Start the server
	r.Run("localhost:8080")
}

// getAllUsersHandler retrieves all users from the repository and returns them in the response.
func getAllUsersHandler(c *gin.Context, userRepo repository.PostgresUserRepository) {
	users, err := userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func createUserHandler(c *gin.Context, userRepo repository.PostgresUserRepository) {
	var newUser repository.User
	var rawRequestBody map[string]interface{}

	// Bind the raw JSON to rawRequestBody and log it
	if err := c.ShouldBindJSON(&rawRequestBody); err != nil {
		fmt.Println("Error binding raw request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println("Received request body:", rawRequestBody)

	// Handle User_id manually if it's present but is an empty string
	if userID, ok := rawRequestBody["User_id"].(string); ok {
		if userID == "" {
			// Set User_id to -1 if it's an empty string
			fmt.Println("User_id is an empty string, setting to -1")
			newUser.User_id = -1
		} else {
			// Try to convert the string to an int
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println("Invalid User_id format, setting to -1:", err)
				newUser.User_id = -1
			} else {
				newUser.User_id = userIDInt
			}
		}
		delete(rawRequestBody, "User_id") // Remove User_id from the map to avoid double unmarshalling
	} else {
		// If User_id is missing, set it to -1
		fmt.Println("User_id is missing, setting to -1")
		newUser.User_id = -1
	}

	// Marshal and Unmarshal the remaining fields to newUser struct
	jsonData, err := json.Marshal(rawRequestBody)
	if err != nil {
		fmt.Println("Error marshalling raw request body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal request body"})
		return
	}

	if err := json.Unmarshal(jsonData, &newUser); err != nil {
		fmt.Println("Error unmarshalling into newUser:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal into user object"})
		return
	}

	// Print the newUser object to ensure it's been populated correctly
	fmt.Printf("User object after unmarshalling: %+v\n", newUser)

	// Create user in the database
	err = userRepo.CreateUser(&newUser)
	if err != nil {
		fmt.Println("Error creating user in the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Debugging info after successful database write
	fmt.Println("User created successfully:", newUser)

	// Respond with the created user
	c.IndentedJSON(http.StatusCreated, newUser)
}

func updateUserHandler(c *gin.Context, userRepo repository.PostgresUserRepository) {
	var updatedUser repository.User
	// Bind the received JSON to updatedUser
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(updatedUser)
	fmt.Println(c.Params)
	// Update the user in the database
	err := userRepo.UpdateUser(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func patchUserHandler(c *gin.Context, userRepo repository.PostgresUserRepository) {
	var updates map[string]interface{}
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam) // Convert the string ID to an integer

	// Bind the received JSON to updates map
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user in the database with only the provided fields
	err = userRepo.PatchUser(userId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUserHandler(c *gin.Context, userRepo repository.PostgresUserRepository) {
	fmt.Println("deleting user")
	// Extract and convert the userId from the URL
	idParam := c.Param("id")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete the user by ID
	err = userRepo.DeleteUserByID(userId)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
