package kademlia

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (node *Node) Cli() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Node CLI started. Commands: put <content>, get <hash>, exit")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.SplitN(input, " ", 2)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "put":
			if len(parts) < 2 {
				fmt.Println("Usage: put <content>")
				continue
			}
			node.Put(parts[1])
		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <hash>")
				continue
			}
			node.Get(parts[1])
		case "exit":
			fmt.Println("Shutting down node.")
			return
		default:
			fmt.Println("Unknown command. Use put <content>, get <hash>, or exit.")
		}
	}
}

func (node *Node) Put(content string) {
	ans, err := node.Client.SendStoreMessage([]byte(content))
	if err != nil {
		fmt.Println("Error storing content:", err)
		return
	}
	fmt.Printf("✅ Content stored!\nHash: %s\nPacket ID: %s\n", ans.Payload.Key, ans.PacketID)
}

func (node *Node) Get(hash string) {
	ans, err := node.Client.SendFindValueMessage(hash)
	if err != nil {
		fmt.Println("Error retrieving content:", err)
		return
	}
	fmt.Printf("✅ Content retrieved!\nHash: %s\nContent: %s\nSource: %s\n", ans.Payload.Key, ans.Payload.Data, ans.Payload.SourceContact.ID.String())
}
