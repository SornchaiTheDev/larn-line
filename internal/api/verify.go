package api

import (
	"github.com/gin-gonic/gin"
)

func VerifyPhoneNumber(c *gin.Context) {

	code := c.Query("code")


	c.JSON(200, gin.H{
		"message": "OK",
		"code":    code,
	})
}
