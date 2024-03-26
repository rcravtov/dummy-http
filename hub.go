package main

import (
	"net"
	"sync"
)

type Connection struct {
	ws   net.Conn
	send chan string
}

type Hub struct {
	connections map[*Connection]struct{}
	connmu      sync.Mutex

	broadcast chan string
}

func NewHub() *Hub {
	hub := &Hub{
		connections: make(map[*Connection]struct{}),
		connmu:      sync.Mutex{},
		broadcast:   make(chan string),
	}

	go func() {
		for {
			m := <-hub.broadcast
			for c := range hub.connections {
				select {
				case c.send <- m:
				default:
					hub.DeleteConnection(c)
				}
			}
		}
	}()

	return hub
}

func (h *Hub) AddConnection(ws net.Conn) *Connection {
	h.connmu.Lock()

	connection := &Connection{
		ws:   ws,
		send: make(chan string),
	}
	h.connections[connection] = struct{}{}

	h.connmu.Unlock()

	return connection
}

func (h *Hub) DeleteConnection(c *Connection) {
	h.connmu.Lock()

	c.ws.Close()
	delete(h.connections, c)
	close(c.send)

	h.connmu.Unlock()
}

func (h *Hub) BroadcastMessage(m string) {
	h.broadcast <- m
}
