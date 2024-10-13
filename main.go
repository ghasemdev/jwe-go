package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"jwe-go/packages/schema"
	"jwe-go/routes"
)

func main() {
	schema.Validate = validator.New(validator.WithRequiredStructEnabled())

	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/encrypt", routes.EncryptEndpoint)
		v1.POST("/decrypt", routes.DecryptEndpoint)
	}

	router.Run(":8080")
}
