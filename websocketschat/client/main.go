package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

func main() {
	// Connect to the WebSocket server
	ws, err := websocket.Dial("ws://localhost:8080/chat", "", "http://localhost/")
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer ws.Close()

	// Read messages from the server in a goroutine
	go func() {
		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				fmt.Println("Error receiving message:", err)
				return
			}
			fmt.Println(msg)
		}
	}()

	// Handle user input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the WebSocket chat!")
	fmt.Println("Commands: join <username>, msg <message>, history, leave")
	fmt.Print("Enter command: ")

	for {
		scanner.Scan()
		input := scanner.Text()

		if strings.HasPrefix(input, "join") {
			// Send the 'join' command to the server
			err := websocket.Message.Send(ws, input)
			if err != nil {
				fmt.Println("Error sending message:", err)
				continue
			}
		} else if strings.HasPrefix(input, "msg") {
			// Send the 'msg' command to the server
			err := websocket.Message.Send(ws, input)
			if err != nil {
				fmt.Println("Error sending message:", err)
				continue
			}
		} else if input == "history" {
			// Send the 'history' command to the server
			err := websocket.Message.Send(ws, input)
			if err != nil {
				fmt.Println("Error sending message:", err)
				continue
			}
		} else if input == "leave" {
			// Send the 'leave' command to the server
			err := websocket.Message.Send(ws, input)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Error sending message:", err)
				continue
			}
			break
		} else {
			fmt.Println("Invalid command. Use: join <username>, msg <message>, leave.")
		}
	}
}
