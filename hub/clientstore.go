package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnStore struct {
	mutex       sync.RWMutex
	connections map[string]map[string]*websocket.Conn
}

// AddConnection adds a connection to the store.
// If the user does not exist in the store, a new entry is created.
// If the user exists but the connection ID is new, the connection is added to the user's connections.
// If the user exists and the connection ID already exists, the connection is updated.
func (cs *ConnStore) AddConnection(userID string, connID string, conn *websocket.Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	if cs.connections[userID] == nil {
		cs.connections[userID] = make(map[string]*websocket.Conn)
	}
	cs.connections[userID][connID] = conn
}

// GetConnection returns a connection for a given user and connection ID.
// If the user does not exist in the store, the boolean return value (exists) is false.
func (cs *ConnStore) GetConnection(userID string, connID string) (*websocket.Conn, bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	if cs.connections[userID] == nil {
		return nil, false
	}
	conn, exists := cs.connections[userID][connID]
	return conn, exists
}

// RemoveConnection removes a connection from the store.
// If the user does not exist in the store, the method does nothing.
// If the user exists but the connection ID does not exist, the method does nothing.
func (cs *ConnStore) RemoveConnection(userID string, connID string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	if cs.connections[userID] == nil {
		return
	}
	delete(cs.connections[userID], connID)
}

// RemoveUser removes a user and all their connections from the store.
// If the user does not exist in the store, the method does nothing.
func (cs *ConnStore) RemoveUser(userID string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.connections, userID)
}

// GetUserConnections returns all the connections for a given user.
// This method is useful for broadcasting messages to all the connections of a user.
// If the user has no connections, an empty slice is returned.
// The returned slice is a copy of the connections, so modifying it will not affect the store.
func (cs *ConnStore) GetUserConnections(userID string) []*websocket.Conn {
	connections := make([]*websocket.Conn, 0)
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	if cs.connections[userID] == nil {
		return connections
	}
	for _, conn := range cs.connections[userID] {
		connections = append(connections, conn)
	}
	return connections
}

// GetAllConnections returns all the connections in the store.
// This method is useful for broadcasting messages to all the connections.
func (cs *ConnStore) GetAllConnections() []*websocket.Conn {
	connections := make([]*websocket.Conn, 0)
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	for _, userConnections := range cs.connections {
		for _, conn := range userConnections {
			connections = append(connections, conn)
		}
	}
	return connections
}

func NewConnStore() *ConnStore {
	return &ConnStore{
		connections: make(map[string]map[string]*websocket.Conn),
	}
}
