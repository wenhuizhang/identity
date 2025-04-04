package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Load and verify an Agent Passport",
}

var loadCmd = &cobra.Command{
	Use:   "load [agent_passport]",
	Short: "Load an Agent Passport",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Loading agent passport")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the loaded Agent Passport",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Validating the loaded Agent Passport")
	},
}

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the loaded Agent Passport",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Forgetting the loaded Agent Passport")
	},
}

func init() {
	VerifyCmd.AddCommand(loadCmd)
	VerifyCmd.AddCommand(validateCmd)
	VerifyCmd.AddCommand(forgetCmd)
}
