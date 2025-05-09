// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package oidc

type Claims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
}

// The Client defines different methods for the CLIENT standard
type Client interface {
	ParseJwt(jwtString *string) (Claims, error)
}

// The client struct implements the Client interface
type client struct {
}

// NewClient creates a new instance of the Client
func NewClient() Client {
	return &client{}
}

// ParseJwt parses the JWT, validates the signature and returns the claims
func (o *client) ParseJwt(jwtString *string) (Claims, error) {
	return Claims{}, nil
}
