package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	Content string
	Sender  string
}

func main() {
	// Connect to the server
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer client.Close()

	fmt.Println("Connected to chat server. ")

	// Get username
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Println("Hello " + username + "! You've joined the chat. Type a message to see the chat history.")

	// Get initial chat history
	var messages []Message
	err = client.Call("ChatService.GetMessages", struct{}{}, &messages)
	if err != nil {
		log.Println("Error getting messages:", err)
	} else {
		printMessages(messages)
	}

	// Main loop
	for {
		fmt.Print("> Enter message (or 'exit' to quit.): ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		text = strings.TrimSpace(text)

		// Check for exit command
		if text == "exit" {
			fmt.Println("Exiting chat...")
			return
		}

		// Send message
		msg := Message{
			Content: text,
			Sender:  username,
		}

		var reply []Message
		err = client.Call("ChatService.SendMessage", msg, &reply)
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}

		// Print updated chat history
		printMessages(reply)
	}
}

// printMessages displays all messages in the chat history
func printMessages(messages []Message) {
	fmt.Println("\n--- Chat History ---")
	for _, msg := range messages {
		fmt.Printf("%s: %s\n", msg.Sender, msg.Content)
	}
	fmt.Println("-------------------\n")
}
