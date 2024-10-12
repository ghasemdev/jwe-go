package routes

import (
	"github.com/gin-gonic/gin"
	"jwe-go/model"
	"jwe-go/utils/json"
	"net/http"
)

func EncryptEndpoint(context *gin.Context) {
	var encryption model.EncryptRequest

	// Use Jsoniter for decoding the JSON body
	body, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Use strict unmarshaling
	if err := utils.StrictUnmarshal(body, &encryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or unknown field"})
		return
	}

	// Manually validate the struct using the validator
	if err := utils.Validate.Struct(encryption); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, encryption)
}
