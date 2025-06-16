package endpoints

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ConnectionManager struct {
	clients map[string]*websocket.Conn
	mux     sync.RWMutex
}

var connMgr = &ConnectionManager{
	clients: make(map[string]*websocket.Conn),
}

func (h *Handler) HandleWS(c *gin.Context) {
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

	id, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	connMgr.mux.Lock()
	connMgr.clients[id] = conn
	connMgr.mux.Unlock()

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			connMgr.mux.Lock()
			delete(connMgr.clients, id)
			connMgr.mux.Unlock()
			return
		}
	}
}
