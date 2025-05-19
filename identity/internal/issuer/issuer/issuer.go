// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/pkg/ptrutil"
)

// getIssuersDirectory returns the path to the issuers directory
func getIssuersDirectory(vaultId string) (string, error) {
	vaultIdDir, err := vault.GetVaultIdDirectory(vaultId)
	if err != nil {
		return "", err
	}

	return filepath.Join(vaultIdDir, "issuers"), nil
}

// GetIssuerIdDirectory returns the path to the issuer ID directory
func GetIssuerIdDirectory(vaultId, issuerId string) (string, error) {
	issuersDir, err := getIssuersDirectory(vaultId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuersDir, issuerId), nil
}

// GetIssuerFilePath returns the path to the issuer file
func GetIssuerFilePath(vaultId, issuerId string) (string, error) {
	issuerIdDir, err := GetIssuerIdDirectory(vaultId, issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuerIdDir, "issuer.json"), nil
}

func saveIssuerConfig(vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig) error {
	// Get the issuer ID directory
	issuerIdDir, err := GetIssuerIdDirectory(vaultId, idpConfig.ClientId)
	if err != nil {
		return err
	}

	// Create the issuer ID directory if it doesn't exist
	if err := os.MkdirAll(issuerIdDir, internalIssuerConstants.DirPerm); err != nil {
		return err
	}

	// Create the issuer config
	issuerConfig := internalIssuerTypes.IssuerConfig{
		IdentityNodeConfig: &internalIssuerTypes.IdentityNodeConfig{
			IdentityNodeAddress: identityNodeAddress,
		},
		IdpConfig: &idpConfig,
	}

	// Marshal the config to JSON
	configData, err := json.Marshal(issuerConfig)
	if err != nil {
		return err
	}

	// Write the config to file
	configFilePath := filepath.Join(issuerIdDir, "idp_config.json")
	if err := os.WriteFile(configFilePath, configData, internalIssuerConstants.FilePerm); err != nil {
		return err
	}

	return nil
}

func getMockIssuerInfo() *string {
	return ptrutil.Ptr("AGNTCY")
}

func RegisterIssuer(
	vaultId, identityNodeAddress string,
	idpConfig internalIssuerTypes.IdpConfig,
) (*coreV1alpha.Issuer, error) {
	// Save the issuer config
	if err := saveIssuerConfig(vaultId, identityNodeAddress, idpConfig); err != nil {
		return nil, err
	}

	// Check connection to identity node
	// Check connection to idp
	// Check if idp is already created locally
	// Check if idp is already registered on the identity node
	// Register idp on the identity node
	issuer := coreV1alpha.Issuer{
		Organization:    getMockIssuerInfo(),
		SubOrganization: getMockIssuerInfo(),
		CommonName:      getMockIssuerInfo(),
	}
	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	registerIssuerRequest := nodeV1alpha.RegisterIssuerRequest{
		Issuer: &issuer,
		Proof:  &proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Registering issuer with request: ", &registerIssuerRequest)

	// Create idp locally in the issuer directory
	issuersDir, err := GetIssuerIdDirectory(vaultId, idpConfig.ClientId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(issuersDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	issuerFilePath, err := GetIssuerFilePath(vaultId, idpConfig.ClientId)
	if err != nil {
		return nil, err
	}

	// Marshal the config to JSON
	issuerData, err := json.Marshal(&issuer)
	if err != nil {
		return nil, err
	}

	// Write the issuer to file
	if err := os.WriteFile(issuerFilePath, issuerData, internalIssuerConstants.FilePerm); err != nil {
		return nil, err
	}

	return &issuer, nil
}

func ListIssuerIds(vaultId string) ([]string, error) {
	// Get the issuers directory
	issuersDir, err := getIssuersDirectory(vaultId)
	if err != nil {
		return nil, err
	}

	// Read the issuers directory
	files, err := os.ReadDir(issuersDir)
	if err != nil {
		return nil, err
	}

	// List the issuer IDs
	var issuerIds []string

	for _, file := range files {
		if file.IsDir() {
			issuerIds = append(issuerIds, file.Name())
		}
	}

	return issuerIds, nil
}

func GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error) {
	// Get the issuer file path
	issuerFilePath, err := GetIssuerFilePath(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	// Read the issuer file
	issuerData, err := os.ReadFile(issuerFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the issuer data
	var issuer coreV1alpha.Issuer
	if err := json.Unmarshal(issuerData, &issuer); err != nil {
		return nil, err
	}

	return &issuer, nil
}

func ForgetIssuer(vaultId, issuerId string) error {
	// Get the issuer directory
	issuerDir, err := GetIssuerIdDirectory(vaultId, issuerId)
	if err != nil {
		return err
	}

	// Check if the issuer directory exists
	if _, err := os.Stat(issuerDir); os.IsNotExist(err) {
		return errors.New("issuer does not exist")
	}

	// Remove the issuer directory
	if err := os.RemoveAll(issuerDir); err != nil {
		return err
	}

	return nil
}

func TestIdpConnection(clientId, clientSecret, issuerUrl string) (*oauth2.Token, error) {
	// Test the connection to the Identity Provider
	ctx := context.Background()

	// Discover OIDC provider config
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		return nil, err
	}

	// Set up the OAuth2 client credentials config
	conf := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     provider.Endpoint().TokenURL,
		Scopes:       []string{},
	}

	// Retrieve a token
	token, err := conf.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}
