package main

import (
	"net"
	"strings"
)

var history = make([][]byte, 0)
var activeUsers = make(map[string]bool)

type WsServer struct {
	clients   map[*Client]bool
	add       chan *Client
	remove    chan *Client
	broadcast chan []byte
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:   make(map[*Client]bool),
		add:       make(chan *Client),
		remove:    make(chan *Client),
		broadcast: make(chan []byte),
	}
}

func (s *WsServer) Run() {
	for {
		select {

		case client := <-s.add:
			s.join(client)
			for _, msg := range history {
				client.send <- msg
			}
		case client := <-s.remove:
			s.quit(client)

		case msg := <-s.broadcast:
			history = append(history, msg)
			activeUsers[strings.Split(string(msg), ":")[0]] = true
			s.broadcastToClients(msg)
		}
	}
}

func (s *WsServer) join(c *Client) {
	s.clients[c] = true
}

func (s *WsServer) quit(c *Client) {
	if _, ok := s.clients[c]; ok {
		delete(s.clients, c)
	}
}

func (s *WsServer) broadcastToClients(message []byte) {
	for c := range s.clients {
		c.send <- message
	}
}

func listenTcp() {
	adr := net.TCPAddr{Port: 8082, IP: net.ParseIP("127.0.0.1")}
	ser, _ := net.ListenTCP("tcp", &adr)
	p:=make([]byte,1024)

	for {
		conn, _ := ser.Accept()
		for key, val := range activeUsers {
			if val == true {
				p = append(p, []byte(key+"\n")...)
			}
		}
		conn.Write(p)
		conn.Close()
	}
}
