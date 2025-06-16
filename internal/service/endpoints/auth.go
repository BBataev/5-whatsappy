package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/BBataev/whatsappy/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h *Handler) HandleLogin(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "invalid request"})
		return
	}

	id, ok, err := postgres.CheckUserCredentials(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := generateToken(id, req.Username, []byte(h.cfg.JWToken))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie(
		"token",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

func (h *Handler) HandleRegister(c *gin.Context) {
	var req RegisterRequest
	id := uuid.New()

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

	if err := postgres.AddNewUser(id, req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	slog.Info("User register", "username", req.Username, "email", req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Register success"})
}

func generateToken(id uuid.UUID, username string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (h *Handler) HandleMe(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no auth token"})
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.cfg.JWToken), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	id, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"username": username,
	})
}
