package endpoints

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	slog.Info("User login...")
	c.JSON(200, gin.H{"message": "Login success"})
}

func HandleRegister(c *gin.Context) {
	slog.Info("User register...")
	c.JSON(200, gin.H{"message": "Register success"})
}
