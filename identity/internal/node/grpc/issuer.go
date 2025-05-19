// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	coreapi "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeapi "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	converters "github.com/agntcy/identity/internal/pkg/converters"
	grpcutil "github.com/agntcy/identity/internal/pkg/grpcutil"
	"github.com/agntcy/identity/pkg/log"
)

type issuerService struct {
	nodeIssuerService node.IssuerService
}

func NewIssuerService(nodeIssuerService node.IssuerService) nodeapi.IssuerServiceServer {
	return &issuerService{
		nodeIssuerService,
	}
}

// Register an issuer by providing the issuer details
func (i *issuerService) Register(
	ctx context.Context,
	req *nodeapi.RegisterIssuerRequest,
) (*nodeapi.RegisterIssuerResponse, error) {
	log.Debug("RegisterIssuer: ", req.Issuer.CommonName)

	// Convert entities and call the node service
	uri, err := i.nodeIssuerService.Register(
		ctx,
		converters.Convert[issuertypes.Issuer](req.Issuer),
		converters.Convert[vctypes.Proof](req.Proof),
	)
	if err != nil {
		return nil, grpcutil.BadRequestError(err)
	}

	// Return the action uri
	return &nodeapi.RegisterIssuerResponse{
		Uri: uri,
	}, nil
}

// Returns the well-known document content for an issuer in
// Json Web Key Set (JWKS) format
func (i *issuerService) GetWellKnown(
	ctx context.Context,
	req *nodeapi.GetIssuerWellKnownRequest,
) (*nodeapi.GetIssuerWellKnownResponse, error) {
	log.Debug("GetIssuerWellKnown: ", req.CommonName)

	// Get the issuer's public keys by common name
	jwks, err := i.nodeIssuerService.GetJwks(ctx, req.CommonName)
	if err != nil {
		return nil, grpcutil.BadRequestError(err)
	}

	return &nodeapi.GetIssuerWellKnownResponse{
		Jwks: converters.Convert[coreapi.Jwks](jwks),
	}, nil
}
