// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"crypto/rand"

	"crypto/x509"

	"encoding/pem"
	"fmt"
	"os"
)

// Note: Temporary key generation function for initial testing purposes. To be removed.

// GenerateECDSAKeyPair generates an ECDSA key pair and encodes them in PEM format.
func GenerateECDSAKeyPair() ([]byte, []byte) {
	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating keys: %v\n", err)
		return nil, nil
	}

	// Encode private key to PEM
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding private key: %v\n", err)
		return nil, nil
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding public key: %v\n", err)
		return nil, nil
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM
}
