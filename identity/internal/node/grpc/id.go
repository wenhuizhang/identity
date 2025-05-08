// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"fmt"

	nodeapi "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
)

type idService struct{}

func NewIdService() nodeapi.IdServiceServer {
	return &idService{}
}

// Generate an Id and its corresponding ResolverMetadata for the specified Issuer
func (idService) Generate(
	ctx context.Context,
	req *nodeapi.GenerateRequest,
) (*nodeapi.GenerateResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// Resolve a specified Id to its corresponding ResolverMetadata
func (idService) Resolve(
	ctx context.Context,
	req *nodeapi.ResolveRequest,
) (*nodeapi.ResolveResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
