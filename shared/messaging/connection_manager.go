package messaging

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/gorilla/websocket"
)

var ErrConnectionNotFound = fmt.Errorf("connection not found")

type connWrapper struct {
	conn *websocket.Conn
	mu sync.Mutex
}

type ConnectionManager struct {
	connections map[string]*connWrapper
	mu sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*connWrapper),
	}
}

func (cm *ConnectionManager) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cm *ConnectionManager) Add(id string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[id] = &connWrapper{
		conn: conn,
		mu: sync.Mutex{},
	}

	log.Printf("Connection added for ID: %s", id)
}

func (cm *ConnectionManager) Remove(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, id)
	log.Printf("Connection removed for ID: %s", id)
}

func (cm *ConnectionManager) Get(id string) (*websocket.Conn, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	wrapper, exists := cm.connections[id]
	if !exists {
		return nil, false
	}
	return wrapper.conn, true
}

func (cm *ConnectionManager) SendMessage(id string, message contracts.WSMessage) error {
	cm.mu.RLock()
	wrapper, exists := cm.connections[id]
	cm.mu.RUnlock()
	if !exists {
		return ErrConnectionNotFound
	}
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	return wrapper.conn.WriteJSON(message)
}