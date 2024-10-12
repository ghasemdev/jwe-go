package main

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/handshake", handshakeEndpoint)
	}

	router.Run(":8080")
}

func handshakeEndpoint(c *gin.Context) {
	response := map[string]string{
		"result": "ok",
	}
	c.JSON(http.StatusOK, response)
}
