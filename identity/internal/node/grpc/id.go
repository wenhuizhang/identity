package grpc

import (
	"context"

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
	return nil, nil
}

// Resolve a specified Id to its corresponding ResolverMetadata
func (idService) Resolve(
	ctx context.Context,
	req *nodeapi.ResolveRequest,
) (*nodeapi.ResolveResponse, error) {
	return nil, nil
}
