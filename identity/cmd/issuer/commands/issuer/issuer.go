// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/spf13/cobra"
)

var issuerFilesystemRepository = filesystem.NewIssuerFilesystemRepository()
var issuerService = issuer.NewIssuerService(issuerFilesystemRepository)

var IssuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "Setup your issuer environment, including your vault, identity provider, and identity network",
	Long: `
The setup command is used to configure your local environment for the Identity CLI tool. With it you can:

- (register) Register with an identity provider, such as DUO or Okta, to manage your Agent and MCP identities
- (list) List your existing issuer configurations
- (show) Show details of an issuer configuration
- (load) Load an issuer configuration
- (forget) Forget an issuer configuration
`,
}

//nolint:mnd // Allow magic number for args
var issuerRegisterCmd = &cobra.Command{
	Use:   "register [identity_node_address] [idp_client_id] [idp_client_secret] [idp_issuer_url]",
	Short: "Register as an Issuer",
	Long:  "Register as an Issuer with an Identity Network using the provided client ID, client secret, and issuer URL.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		identityNodeAddress := args[0]
		clientID := args[1]
		clientSecret := args[2]
		issuerURL := args[3]

		config := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		issuerId, err := issuerService.RegisterIssuer(cache.VaultId, identityNodeAddress, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error registering as an Issuer: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully registered as an Issuer with ID: %s\n", issuerId)

		// Update the cache with the new issuer ID
		cache.IssuerId = issuerId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
	},
}

var issuerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing issuer configurations",
	Long:  "List your existing issuer configurations",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		issuers, err := issuerService.GetAllIssuers(cache.VaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing issuers: %v\n", err)
			return
		}
		if len(issuers) == 0 {
			fmt.Fprintf(os.Stdout, "No issuers found.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Existing issuers:\n")
		for _, issuer := range issuers {
			fmt.Fprintf(os.Stdout, "- %s, %s\n", issuer.Id, *issuer.Issuer.CommonName)
		}
	},
}
var issuerShowCmd = &cobra.Command{
	Use:   "show [issuer_id]",
	Short: "Show details of an issuer configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		issuerId := args[0]

		issuer, err := issuerService.GetIssuer(cache.VaultId, issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stdout, "No issuer found with ID: %s\n", issuerId)
			return
		}

		issuerJSON, err := json.MarshalIndent(issuer, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(issuerJSON))
	},
}

var issuerForgetCmd = &cobra.Command{
	Use:   "forget [issuer_id]",
	Short: "Forget an issuer configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		issuerId := args[0]

		err = issuerService.ForgetIssuer(cache.VaultId, issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting issuer: %v\n", err)
			return
		}

		// If the issuer was the current issuer in the cache, clear the cache of issuer, metadata, and badge IDs
		if cache.IssuerId == issuerId {
			cache.IssuerId = ""
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot issuer with ID: %s\n", issuerId)
	},
}

var issuerLoadCmd = &cobra.Command{
	Use:   "load [issuer_id]",
	Short: "Load an issuer configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		issuerId := args[0]

		// check the issuer id is valid
		issuer, err := issuerService.GetIssuer(cache.VaultId, issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stderr, "No issuer found with ID: %s\n", issuerId)
			return
		}

		// save the issuer id to the cache
		cache.IssuerId = issuerId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded issuer with ID: %s\n", issuerId)

	},
}

func init() {
	IssuerCmd.AddCommand(issuerRegisterCmd)
	IssuerCmd.AddCommand(issuerListCmd)
	IssuerCmd.AddCommand(issuerShowCmd)
	IssuerCmd.AddCommand(issuerForgetCmd)
	IssuerCmd.AddCommand(issuerLoadCmd)
}
