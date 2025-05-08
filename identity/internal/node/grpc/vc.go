package grpc

import (
	"context"

	nodeapi "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type vcService struct{}

func NewVcService() nodeapi.VcServiceServer {
	return &vcService{}
}

// Publish an issued Verifiable Credential
func (vcService) Publish(
	ctx context.Context,
	req *nodeapi.PublishRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

// Verify an existing Verifiable Credential
func (vcService) Verify(
	ctx context.Context,
	req *nodeapi.VerifyRequest,
) (*emptypb.Empty, error) {
	return nil, nil
}

// Returns the well-known Verifiable Credentials for the specified Id
func (vcService) GetWellKnown(
	ctx context.Context,
	req *nodeapi.GetVcWellKnownRequest,
) (*nodeapi.GetVcWellKnownResponse, error) {
	return nil, nil
}

// Search for Verifiable Credentials based on the specified criteria
func (vcService) Search(
	ctx context.Context,
	req *nodeapi.SearchRequest,
) (*nodeapi.SearchResponse, error) {
	return nil, nil
}
