package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[int]*websocket.Conn)
	clientsMu sync.Mutex
)

func AddWebSocketConn(userID int, conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[userID] = conn
}

func RemoveWebSocketConn(userID int) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	if _, ok := clients[userID]; ok {
		delete(clients, userID)
	}
}

func GetWebSocketConn(userID int) *websocket.Conn {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	return clients[userID]
}
