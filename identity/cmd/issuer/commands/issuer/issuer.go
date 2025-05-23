// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

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

var (
	// setup the issuer service
	issuerFilesystemRepository = filesystem.NewIssuerFilesystemRepository()
	issuerService              = issuer.NewIssuerService(issuerFilesystemRepository)

	// setup the command flags
	registerIdentityNodeAddress string
	registerClientID            string
	registerClientSecret        string
	registerIssuerURL           string
	registerOrganization        string
	registerSubOrganization     string
	showIssuerId                string
	forgetIssuerId              string
	loadIssuerId                string
)

var IssuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "Register as an Issuer and manage issuer configurations",
	Long: `
The issuer command is used to register as an Issuer and manage issuer configurations.
`,
}

//nolint:lll // Allow long lines for CLI
var issuerRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register as an Issuer",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForIssuer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the identity node address is not set, prompt the user for it interactively
		if registerIdentityNodeAddress == "" {
			fmt.Fprintf(os.Stdout, "Identity node address (default %s): ", defaultNodeAddress)
			_, err = fmt.Scanln(&registerIdentityNodeAddress)
			if err != nil {
				// If the user just presses Enter, registerIdentityNodeAddress will be "" and err will be an "unexpected newline" error.
				// We should allow this and use the default value.
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading identity node address: %v\n", err)
					return
				}
			}
		}
		// If no address was entered (input was empty or only whitespace), use the default.
		if registerIdentityNodeAddress == "" {
			registerIdentityNodeAddress = defaultNodeAddress
		}

		// if the client ID is not set, prompt the user for it interactively
		if registerClientID == "" {
			fmt.Fprintf(os.Stdout, "IdP client ID: ")
			_, err = fmt.Scanln(&registerClientID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP client ID: %v\n", err)
				return
			}
			if registerClientID == "" {
				fmt.Fprintf(os.Stderr, "IdP Client ID cannot be empty.\n")
				return
			}
		}

		// if the client secret is not set, prompt the user for it interactively
		if registerClientSecret == "" {
			fmt.Fprintf(os.Stdout, "IdP client secret: ")
			_, err = fmt.Scanln(&registerClientSecret)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP client secret: %v\n", err)
				return
			}
			if registerClientSecret == "" {
				fmt.Fprintf(os.Stderr, "IdP Client secret cannot be empty.\n")
				return
			}
		}

		// if the issuer URL is not set, prompt the user for it interactively
		if registerIssuerURL == "" {
			fmt.Fprintf(os.Stdout, "IdP issuer URL: ")
			_, err = fmt.Scanln(&registerIssuerURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP issuer URL: %v\n", err)
				return
			}
			if registerIssuerURL == "" {
				fmt.Fprintf(os.Stderr, "IdP issuer URL cannot be empty.\n")
				return
			}
		}

		// if the organization is not set, prompt the user for it interactively
		if registerOrganization == "" {
			fmt.Fprintf(os.Stdout, "Organization name: ")
			_, err = fmt.Scanln(&registerOrganization)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading organization name: %v\n", err)
				return
			}
			if registerOrganization == "" {
				fmt.Fprintf(os.Stderr, "Organization name cannot be empty.\n")
				return
			}
		}

		// if the sub-organization is not set, prompt the user for it interactively
		if registerSubOrganization == "" {
			fmt.Fprintf(os.Stdout, "Sub-organization name (default %s): ", registerOrganization)
			_, err = fmt.Scanln(&registerSubOrganization)
			if err != nil {
				// If the user just presses Enter, registerSubOrganization will be "" and err will be an "unexpected newline" error.
				// We should allow this and use the default value.
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading sub-organization name: %v\n", err)
					return
				}
				if registerSubOrganization == "" {
					registerSubOrganization = registerOrganization
				}
			}
		}

		idpConfig := issuerTypes.IdpConfig{
			ClientId:     registerClientID,
			ClientSecret: registerClientSecret,
			IssuerUrl:    registerIssuerURL,
		}

		// extract the root url from the issuer URL as the common name
		issuerUrl, err := url.Parse(idpConfig.IssuerUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing issuer URL: %v\n", err)
			return
		}
		commonName := issuerUrl.Hostname()
		if commonName == "" {
			fmt.Fprintf(os.Stderr, "Error extracting common name from issuer URL: %v\n", err)
			return
		}

		coreIssuer := coreV1alpha.Issuer{
			Organization:    &registerOrganization,
			SubOrganization: &registerSubOrganization,
			CommonName:      &commonName,
		}

		identityNodeConfig := issuerTypes.IdentityNodeConfig{
			IdentityNodeAddress: registerIdentityNodeAddress,
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
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}

//nolint:lll // Allow long lines for CLI
var issuerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing issuer configurations",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForIssuer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
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

//nolint:lll // Allow long lines for CLI
var issuerShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of an issuer configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		// if the issuer id is not set, prompt the user for it interactively
		if showIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to show:\n")
			_, err := fmt.Scanln(&showIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if showIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		issuer, err := issuerService.GetIssuer(cache.VaultId, showIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stdout, "No issuer found with ID: %s\n", showIssuerId)
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

//nolint:lll // Allow long lines for CLI
var issuerForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget an issuer configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForIssuer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the issuer id is not set, prompt the user for it interactively
		if forgetIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to forget:\n")
			_, err := fmt.Scanln(&forgetIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if forgetIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		err = issuerService.ForgetIssuer(cache.VaultId, forgetIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting issuer: %v\n", err)
			return
		}

		// If the issuer was the current issuer in the cache, clear the cache of issuer, metadata, and badge IDs
		if cache.IssuerId == forgetIssuerId {
			cache.IssuerId = ""
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot issuer with ID: %s\n", forgetIssuerId)
	},
}

//nolint:lll // Allow long lines for CLI
var issuerLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load an issuer configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForIssuer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the issuer id is not set, prompt the user for it interactively
		if loadIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to load:\n")
			_, err := fmt.Scanln(&loadIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if loadIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		// check the issuer id is valid
		issuer, err := issuerService.GetIssuer(cache.VaultId, loadIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stderr, "No issuer found with ID: %s\n", loadIssuerId)
			return
		}

		// save the issuer id to the cache
		cache.IssuerId = loadIssuerId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded issuer with ID: %s\n", loadIssuerId)

	},
}

//nolint:lll // Allow long lines for CLI
func init() {
	issuerRegisterCmd.Flags().StringVarP(&registerIdentityNodeAddress, "identity-node-address", "i", "", "Identity node address")
	issuerRegisterCmd.Flags().StringVarP(&registerClientID, "client-id", "c", "", "IdP client ID")
	issuerRegisterCmd.Flags().StringVarP(&registerClientSecret, "client-secret", "s", "", "IdP client secret")
	issuerRegisterCmd.Flags().StringVarP(&registerIssuerURL, "issuer-url", "u", "", "IdP issuer URL")
	issuerRegisterCmd.Flags().StringVarP(&registerOrganization, "organization", "o", "", "Organization name")
	issuerRegisterCmd.Flags().StringVarP(&registerSubOrganization, "sub-organization", "b", "", "Sub-organization name")
	IssuerCmd.AddCommand(issuerRegisterCmd)

	IssuerCmd.AddCommand(issuerListCmd)

	issuerShowCmd.Flags().StringVarP(&showIssuerId, "issuer-id", "i", "", "The ID of the issuer to show")
	IssuerCmd.AddCommand(issuerShowCmd)

	issuerForgetCmd.Flags().StringVarP(&forgetIssuerId, "issuer-id", "i", "", "The ID of the issuer to forget")
	IssuerCmd.AddCommand(issuerForgetCmd)

	issuerLoadCmd.Flags().StringVarP(&loadIssuerId, "issuer-id", "i", "", "The ID of the issuer to load")
	IssuerCmd.AddCommand(issuerLoadCmd)
}
