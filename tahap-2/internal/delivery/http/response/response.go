package response

import (
	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Status  string      `json:"status,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(c *gin.Context, statusCode int, result interface{}) {
	c.JSON(statusCode, JSONResponse{
		Status: "SUCCESS",
		Result: result,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, JSONResponse{
		Message: message,
	})
}
