package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage your connection to an Identity Network node",
}

var networkConnectCmd = &cobra.Command{
	Use:   "setup [identity_node_address]",
	Short: "Setup the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up connection to an Identity Network node")
	},
}

var networkTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing connection to an Identity Network node")
	},
}

var networkForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Forgetting connection to an Identity Network node")
	},
}

func init() {
	NetworkCmd.AddCommand(networkConnectCmd)
	NetworkCmd.AddCommand(networkTestCmd)
	NetworkCmd.AddCommand(networkForgetCmd)
}
