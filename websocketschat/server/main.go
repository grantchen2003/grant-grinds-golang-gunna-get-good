package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn *websocket.Conn
	name string
}

type Server struct {
	broadcastChan   chan string
	chatHistory     []string
	chatHistorySize int
	clients         map[*Client]bool
	mutex           sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		broadcastChan:   make(chan string),
		chatHistory:     []string{},
		chatHistorySize: 100,
		clients:         make(map[*Client]bool),
		mutex:           sync.RWMutex{},
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	client := &Client{conn: ws}

	log.Println("New incoming connection from client:", client.conn.Request().RemoteAddr)

	s.mutex.Lock()
	s.clients[client] = true
	s.mutex.Unlock()

	defer s.deleteClient(client)

	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				s.deleteClient(client)
				log.Println("Connection ended from client:", client.conn.RemoteAddr())
				return
			}
			log.Println("Error reading message:", err)
			return
		}
		s.handleCommand(client, msg)
	}
}

func (s *Server) handleCommand(client *Client, msg string) {
	parts := strings.SplitN(msg, " ", 2)
	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	switch cmd {
	case "join":
		if client.name != "" {
			client.conn.Write([]byte("You have already joined the chat\n"))
			return
		}
		client.name = arg
		s.broadcastChan <- fmt.Sprintf("%s joined the chat", client.name)
	case "msg":
		if client.name == "" {
			client.conn.Write([]byte("You must join first\n"))
			return
		}
		s.broadcastChan <- fmt.Sprintf("%s: %s", client.name, arg)
	case "history":
		s.mutex.Lock()
		client.conn.Write([]byte(strings.Join(s.chatHistory, "\n")))
		s.mutex.Unlock()
	case "leave":
		s.deleteClient(client)
		s.broadcastChan <- fmt.Sprintf("%s left the chat", client.name)
	}
}

func (s *Server) deleteClient(client *Client) {
	s.mutex.Lock()
	delete(s.clients, client)
	client.conn.Close()
	s.mutex.Unlock()
}

func (s *Server) Broadcast() {
	for msg := range s.broadcastChan {
		s.mutex.Lock()
		for client := range s.clients {
			err := websocket.Message.Send(client.conn, fmt.Sprintf("Server: %s", msg))
			if err != nil {
				log.Println(err)
			}
		}
		s.chatHistory = append(s.chatHistory, msg)
		if len(s.chatHistory) > s.chatHistorySize {
			s.chatHistory = s.chatHistory[1:]
		}
		log.Printf("Broadcasted message: %s", msg)
		s.mutex.Unlock()
	}
}

func main() {
	server := NewServer()

	go server.Broadcast()

	http.Handle("/chat", websocket.Handler(server.HandleWS))

	fmt.Println("Starting websockets chat server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
