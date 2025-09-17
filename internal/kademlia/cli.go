package kademlia

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (node *Node) Cli() {
	// Root command
	var rootCmd = &cobra.Command{
		Use:   "node",
		Short: "Kademlia node CLI",
		Long:  "A command-line interface for interacting with a running Kademlia node.",
	}

	// put command
	var putCmd = &cobra.Command{
		Use:   "put [content]",
		Short: "Store content in the DHT",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node.Put(args[0])
		},
	}

	// get command
	var getCmd = &cobra.Command{
		Use:   "get [hash]",
		Short: "Retrieve content from the DHT by hash",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node.Get(args[0])
		},
	}

	// exit command
	var exitCmd = &cobra.Command{
		Use:   "exit",
		Short: "Shut down the node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Shutting down node.")
			os.Exit(0)
		},
	}

	// Attach subcommands to root
	rootCmd.AddCommand(putCmd, getCmd, exitCmd)

	// Run CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
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
	fmt.Printf("✅ Content retrieved!\nHash: %s\nContent: %s\nSource: %s\n", ans.Payload.Key, ans.Payload.Data, ans.Payload.SourceContact)
}
