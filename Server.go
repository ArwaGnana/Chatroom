package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// Message represents a chat message
type Message struct {
	Content string
	Sender  string
}

// ChatService represents the RPC service
type ChatService struct {
	messages []Message
	mu       sync.Mutex
}

// NewChatService creates a new ChatService
func NewChatService() *ChatService {
	return &ChatService{
		messages: make([]Message, 0),
	}
}

// SendMessage adds a new message to the chat history
func (cs *ChatService) SendMessage(msg Message, reply *[]Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if msg.Content == "" {
		return errors.New("message content cannot be empty")
	}

	cs.messages = append(cs.messages, msg)
	*reply = cs.messages

	// Display only the new message on the server
	fmt.Printf("[NEW MESSAGE] %s: %s\n", msg.Sender, msg.Content)

	return nil
}

// GetMessages returns the chat history
func (cs *ChatService) GetMessages(_ struct{}, reply *[]Message) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	*reply = cs.messages
	return nil
}

func main() {
	chatService := NewChatService()

	// Register the service
	err := rpc.Register(chatService)
	if err != nil {
		log.Fatal("Error registering service:", err)
	}

	// Set up TCP listener
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	fmt.Println("Chat server started on port 1234")
	fmt.Println("Waiting for messages...")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a goroutine
		go rpc.ServeConn(conn)
	}
}
