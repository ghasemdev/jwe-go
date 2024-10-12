package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"jwe-go/packages/json"
	"jwe-go/packages/schema"
	"jwe-go/routes"
	"net/http"
)

type User struct {
	Name  string `json:"name" validate:"required,min=3"`  // Name is required and should be at least 3 characters
	Email string `json:"email" validate:"required,email"` // Email is required and should be a valid email format
}

func main() {
	schema.Validate = validator.New(validator.WithRequiredStructEnabled())

	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/user", getUserEndpoint)
		v1.POST("/encrypt", routes.EncryptEndpoint)
	}

	router.Run(":8080")
}

func getUserEndpoint(c *gin.Context) {
	// GET request: Return a JSON response using Jsoniter
	user := User{Name: "John Doe", Email: "john.doe@example.com"}

	// Use Jsoniter for marshaling
	jsonData, err := json.CONFIG.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON encoding failed"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}
