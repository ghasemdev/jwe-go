package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	// "github.com/go-jose/go-jose/v4"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type User struct {
	Name  string `json:"name" validate:"required,min=3"`  // Name is required and should be at least 3 characters
	Email string `json:"email" validate:"required,email"` // Email is required and should be a valid email format
}

type Encryption struct {
	Plaintext  string `json:"plaintext" validate:"required"`
	PublicKey string `json:"publicKey" validate:"required"`
}

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/handshake", handshakeEndpoint)
		v1.GET("/user", getUserEndpoint)
		v1.POST("/user", postUserEndpoint)
		v1.POST("/encrypt", encryptEndpoint)
	}

	router.Run(":8080")
}

func handshakeEndpoint(c *gin.Context) {
	response := map[string]string{
		"result": "ok",
	}
	c.JSON(http.StatusOK, response)
}

func getUserEndpoint(c *gin.Context) {
	// GET request: Return a JSON response using Jsoniter
	user := User{Name: "John Doe", Email: "john.doe@example.com"}

	// Use Jsoniter for marshaling
	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON encoding failed"})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}

func postUserEndpoint(c *gin.Context) {
	var user User

	// Use Jsoniter for decoding the JSON body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Use strict unmarshaling
	if err := strictUnmarshal(body, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or unknown field"})
		return
	}

	// Manually validate the struct using the validator
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func encryptEndpoint(c *gin.Context) {
	var encryption Encryption

	// Use Jsoniter for decoding the JSON body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Use strict unmarshaling
	if err := strictUnmarshal(body, &encryption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or unknown field"})
		return
	}

	// Manually validate the struct using the validator
	if err := validate.Struct(encryption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, encryption)
}

// Custom unmarshaler that rejects extra fields
func strictUnmarshal(data []byte, v interface{}) error {
	// Unmarshal into the actual struct
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields() // Reject unknown fields
	return decoder.Decode(v)
}
