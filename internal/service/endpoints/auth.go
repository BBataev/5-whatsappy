package endpoints

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/BBataev/whatsappy/internal/config"
	"github.com/BBataev/whatsappy/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func HandleLogin(c *gin.Context, cfg *config.Config) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "invalid request"})
		return
	}

	ok, err := postgres.CheckUserCredentials(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := generateToken(req.Username, []byte(cfg.JWToken))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func HandleRegister(c *gin.Context) {
	var req RegisterRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("Failed to parse request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	conflict, err := postgres.CheckUserConflict(req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if conflict {
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already in use"})
		return
	}

	if err := postgres.AddNewUser(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert user"})
		return
	}

	slog.Info("User register", "username", req.Username, "email", req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Register success"})
}

func generateToken(username string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
