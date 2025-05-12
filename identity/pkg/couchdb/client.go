// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package couchdb

import (
	"context"
	"fmt"
	"net"

	"github.com/agntcy/identity/pkg/log"
	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

// Connect : This is a helper function to connect couchDB
func Connect(
	ctx context.Context,
	url, port, username, password string,
) (*kivik.Client, error) {
	host := net.JoinHostPort(url, port)

	// Set connection string
	connectionString := fmt.Sprintf("http://%s", host)

	// Check if auth is enabled
	if (username != "") && (password != "") {
		connectionString = fmt.Sprintf(
			"http://%s:%s@%s",
			username,
			password,
			host,
		)
	}

	log.Debug("Connecting to CouchDB:", connectionString)

	client, err := kivik.New("couch", connectionString)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Disconnect(_ context.Context, client *kivik.Client) error {
	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
