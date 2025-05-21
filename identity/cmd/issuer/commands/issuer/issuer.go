// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/spf13/cobra"
)

const (
	defaultNodeAddress = "http://localhost:4000"
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
	Use:   "register",
	Short: "Register as an Issuer",
	Long:  "Register as an Issuer with an Identity Network using the provided client ID, client secret, and issuer URL.",
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

		// prompt user for identity node address, default to localhost:4000 if not provided
		fmt.Fprintf(os.Stdout, "Enter the identity node address (default %s: ", defaultNodeAddress)
		var identityNodeAddress string
		_, err = fmt.Scanln(&identityNodeAddress)
		if err != nil {
			// If the user just presses Enter, identityNodeAddress will be "" and err will be an "unexpected newline" error.
			// We should allow this and use the default value.
			if err.Error() != "unexpected newline" {
				fmt.Fprintf(os.Stderr, "Error reading identity node address: %v\n", err)
				return
			}
		}

		// If no address was entered (input was empty or only whitespace), use the default.
		if identityNodeAddress == "" {
			identityNodeAddress = defaultNodeAddress
		}

		// prompt user for client ID, client secret, and issuer URL
		fmt.Fprintf(os.Stdout, "Enter the IdP client ID: ")
		var clientID string
		_, err = fmt.Scanln(&clientID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading IdP client ID: %v\n", err)
			return
		}
		if clientID == "" {
			fmt.Fprintf(os.Stderr, "IdP Client ID cannot be empty.\n")
			return
		}

		fmt.Fprintf(os.Stdout, "Enter the IdP client secret: ")
		var clientSecret string
		_, err = fmt.Scanln(&clientSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading IdP client secret: %v\n", err)
			return
		}
		if clientSecret == "" {
			fmt.Fprintf(os.Stderr, "IdP Client secret cannot be empty.\n")
			return
		}

		fmt.Fprintf(os.Stdout, "Enter the IdP issuer URL: ")
		var issuerURL string
		_, err = fmt.Scanln(&issuerURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading IdP issuer URL: %v\n", err)
			return
		}
		if issuerURL == "" {
			fmt.Fprintf(os.Stderr, "IdP issuer URL cannot be empty.\n")
			return
		}

		idpConfig := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		// extract the root url from the issuer URL as the common name
		commonNamePattern := `^https?://([^/]+)`
		re := regexp.MustCompile(commonNamePattern)
		commonNameMatches := re.FindStringSubmatch(idpConfig.IssuerUrl)
		if len(commonNameMatches) < 2 {
			fmt.Fprintf(os.Stderr, "Error extracting common name from issuer URL: %s\n", idpConfig.IssuerUrl)
			return
		}
		commonName := commonNameMatches[1]

		// prompt user for organization and sub-organization
		fmt.Fprintf(os.Stdout, "Enter your organization name: ")
		var organization string
		_, err = fmt.Scanln(&organization)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading organization name: %v\n", err)
			return
		}
		if organization == "" {
			fmt.Fprintf(os.Stderr, "Organization name cannot be empty.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Enter your sub-organization name (default %s): ", organization)
		var subOrganization string
		_, err = fmt.Scanln(&subOrganization)
		if err != nil {
			// If the user just presses Enter, subOrganization will be "" and err will be an "unexpected newline" error.
			// We should allow this and use the organization default value.
			if err.Error() != "unexpected newline" {
				fmt.Fprintf(os.Stderr, "Error reading sub-organization name: %v\n", err)
				return
			}
		}
		// If no sub-organization was entered (input was empty or only whitespace), use the organization as the default.
		if subOrganization == "" {
			subOrganization = organization
		}

		coreIssuer := coreV1alpha.Issuer{
			Organization:    &organization,
			SubOrganization: &subOrganization,
			CommonName:      &commonName,
		}

		identityNodeConfig := issuerTypes.IdentityNodeConfig{
			IdentityNodeAddress: identityNodeAddress,
		}

		issuer := issuerTypes.Issuer{
			Id:                 idpConfig.ClientId,
			Issuer:             &coreIssuer,
			IdentityNodeConfig: &identityNodeConfig,
			IdpConfig:          &idpConfig,
		}

		issuerId, err := issuerService.RegisterIssuer(cache.VaultId, &issuer)
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
