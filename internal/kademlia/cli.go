package kademlia

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (node *Node) Cli(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	fmt.Fprintln(out, "Node CLI started. Commands: put <content>, get <hash>, exit")

	for {
		fmt.Fprintln(out, "Commands: put <content>, get <hash>, exit")
		fmt.Fprint(out, "> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.SplitN(input, " ", 2)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "put":
			if len(parts) < 2 {
				fmt.Fprintln(out, "Usage: put <content>")
				continue
			}
			result, err := node.Put(parts[1])
			if err != nil {
				fmt.Fprintln(out, "Error storing content:", err)
			} else {
				fmt.Fprint(out, result)
			}
		case "get":
			if len(parts) < 2 {
				fmt.Fprintln(out, "Usage: get <hash>")
				continue
			}
			result, err := node.Get(parts[1])
			if err != nil {
				fmt.Fprintln(out, "Error retrieving content:", err)
			} else {
				fmt.Fprint(out, result)
			}
		case "exit":
			fmt.Fprintln(out, "Shutting down node.")
			return
		case "print":
			node.PrintRoutingTable()
		default:
			fmt.Fprintln(out, "Unknown command. Use put <content>, get <hash>, or exit.")
		}
	}
}

func (node *Node) Put(content string) (string, error) {
	ans, err := node.Client.SendStoreMessage([]byte(content))
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("Content stored!\nHash: %s\nPacket ID: %s\n", ans.Payload.Key, ans.PacketID)
	return result, nil
}

func (node *Node) Get(hash string) (string, error) {
	ans, err := node.Client.SendFindValueMessage(hash)
	if err != nil {
		return "", err
	}
	if ans.Payload.SourceContact.ID == nil {
		return "", fmt.Errorf("no source contact found")
	}
	result := fmt.Sprintf("Content retrieved!\nHash: %s\nContent: %s\nSource: %s\n", ans.Payload.Key, ans.Payload.Data, ans.Payload.SourceContact.ID.String())
	return result, nil
}
