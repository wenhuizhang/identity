package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Create and manage your agent artefact versions",
}

var versionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your existing agent artefact versions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all of your existing agent artefacts")
	},
}

var versionLoadCmd = &cobra.Command{
	Use:   "load [version_id]",
	Short: "Load an existing agent artefact version <version_id>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Loading agent artefact version %s\n", args[0])
	},
}

var versionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently loaded agent artefact version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Showing the currently loaded agent artefact version")
	},
}

var versionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and loads a new agent artefact version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating a new agent artefact version")
	},
}

var versionForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current agent artefact version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Forgetting the current agent artefact version")
	},
}

var versionPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current agent artefact version identity",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Publishing the current agent artefact version identity")
	},
}

func init() {
	versionCmd.AddCommand(versionListCmd)
	versionCmd.AddCommand(versionLoadCmd)
	versionCmd.AddCommand(versionShowCmd)
	versionCmd.AddCommand(versionCreateCmd)
	versionCmd.AddCommand(versionForgetCmd)
	versionCmd.AddCommand(versionPublishCmd)
}
