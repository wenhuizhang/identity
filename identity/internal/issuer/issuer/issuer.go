// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package issuer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/agntcy/identity/internal/core/issuer/types"
	vcTypes "github.com/agntcy/identity/internal/core/vc/types"
	issuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	issuerTypesInternal "github.com/agntcy/identity/internal/issuer/types"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type registerIssuerRequest struct {
	Issuer types.Issuer  `json:"issuer"`
	Proof  vcTypes.Proof `json:"proof"`
}

// getIssuersDirectory returns the path to the issuers directory
func getIssuersDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".identity", "issuers"), nil
}

// GetIssuerIdDirectory returns the path to the issuer ID directory
func GetIssuerIdDirectory(issuerId string) (string, error) {
	issuersDir, err := getIssuersDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(issuersDir, issuerId), nil
}

// GetIssuerFilePath returns the path to the issuer file
func GetIssuerFilePath(issuerId string) (string, error) {
	issuerIdDir, err := GetIssuerIdDirectory(issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuerIdDir, "issuer.json"), nil
}

func RegisterIssuer(identityNodeAddress string, idpConfig issuerTypesInternal.IdpConfig) (*types.Issuer, error) {
	// Check connection to identity node
	// Check connection to idp
	// Check if idp is already created locally
	// Check if idp is already registered on the identity node
	// Register idp on the identity node
	issuer := types.Issuer{
		Organization:    "AGNTCY",
		SubOrganization: "AGNTCY",
		CommonName:      "AGNTCY",
	}
	proof := vcTypes.Proof{
		Type:         "RsaSignature2018",
		ProofPurpose: "assertionMethod",
		ProofValue:   "example-proof-value",
	}

	registerIssuerRequest := registerIssuerRequest{
		Issuer: issuer,
		Proof:  proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Registering issuer with request: ", registerIssuerRequest)

	// Create idp locally in the issuer directory
	issuersDir, err := GetIssuerIdDirectory(idpConfig.ClientId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(issuersDir, issuerConstants.DirPerm); err != nil {
		return nil, err
	}

	issuerFilePath, err := GetIssuerFilePath(idpConfig.ClientId)
	if err != nil {
		return nil, err
	}

	// Marshal the config to JSON
	issuerData, err := json.Marshal(issuer)
	if err != nil {
		return nil, err
	}

	// Write the config to file
	if err := os.WriteFile(issuerFilePath, issuerData, issuerConstants.FilePerm); err != nil {
		return nil, err
	}

	return &issuer, nil
}

func ListIssuerIds() ([]string, error) {
	// Get the issuers directory
	issuersDir, err := getIssuersDirectory()
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

func GetIssuer(issuerId string) (*types.Issuer, error) {
	// Get the issuer file path
	issuerFilePath, err := GetIssuerFilePath(issuerId)
	if err != nil {
		return nil, err
	}

	// Read the issuer file
	issuerData, err := os.ReadFile(issuerFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the issuer data
	var issuer types.Issuer
	if err := json.Unmarshal(issuerData, &issuer); err != nil {
		return nil, err
	}

	return &issuer, nil
}

func ForgetIssuer(issuerId string) error {
	// Get the issuer directory
	issuerDir, err := GetIssuerIdDirectory(issuerId)
	if err != nil {
		return err
	}

	// Check if the issuer directory exists
	if _, err := os.Stat(issuerDir); os.IsNotExist(err) {
		return errors.New("Issuer does not exist")
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
