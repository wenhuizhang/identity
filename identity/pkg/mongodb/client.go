// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package mongodb

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/agntcy/identity/pkg/log"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// Connect : This is a helper function to connect mongoDB
func Connect(
	ctx context.Context,
	url, port, username, password string,
	enableTracing bool,
	enableDbLogs bool,
) (*mongo.Client, error) {
	host := net.JoinHostPort(url, port)

	// Set connection string
	connectionString := fmt.Sprintf("mongodb://%s", host)

	// Check if auth is enabled
	if (username != "") && (password != "") {
		// Certificate path
		certPath, _ := filepath.Abs("./global-bundle.pem")

		connectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s%s",
			username,
			password,
			host,
			fmt.Sprintf(
				"/?tls=true&tlsCAFile=%s&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false",
				certPath,
			),
		)
	}

	// Set client options
	clientOptions := options.Client().
		ApplyURI(connectionString)

	if enableDbLogs {
		// Add log monitor if debug mode
		logMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				log.Debug("Command started: ", evt.Command)
			},
		}
		clientOptions = clientOptions.SetMonitor(logMonitor)
	}

	if enableTracing {
		clientOptions = clientOptions.SetMonitor(otelmongo.NewMonitor())
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Disconnect(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
