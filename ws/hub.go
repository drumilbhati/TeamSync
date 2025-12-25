package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// teamID -> list of connections
	teams map[int]map[*websocket.Conn]bool

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		teams: make(map[int]map[*websocket.Conn]bool),
	}
}

func (h *Hub) AddUser(conn *websocket.Conn, teamIDs []int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, teamID := range teamIDs {
		if _, ok := h.teams[teamID]; !ok {
			h.teams[teamID] = make(map[*websocket.Conn]bool)
		}
		h.teams[teamID][conn] = true
	}
}

func (h *Hub) RemoveUser(conn *websocket.Conn, teamIDs []int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, teamID := range teamIDs {
		if members, ok := h.teams[teamID]; ok {
			delete(members, conn)
			if len(members) == 0 {
				delete(h.teams, teamID)
			}
		}
	}
	conn.Close()
}

func (h *Hub) BroadcastToTeam(teamID int, message []byte) {
	h.mu.Lock()

	var connections []*websocket.Conn
	for conn := range h.teams[teamID] {
		connections = append(connections, conn)
	}
	h.mu.Unlock()

	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			conn.Close()
		}
	}
}
