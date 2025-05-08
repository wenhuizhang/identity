package grpc

import (
	"context"

	issuerapi "github.com/agntcy/identity/api/agntcy/identity/issuer/v1alpha1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type localService struct{}

func NewLocalService() issuerapi.LocalServiceServer {
	return &localService{}
}

// Generate a keypair in Json Web Key (JWK) format
func (localService) KeyGen(
	ctx context.Context,
	req *emptypb.Empty,
) (*issuerapi.KeyGenResponse, error) {
	return nil, nil
}

// Issue a Verifiable Credential in a specific Envelope Type
func (localService) IssueVC(
	ctx context.Context,
	req *issuerapi.IssueVCRequest,
) (*issuerapi.IssueVCResponse, error) {
	return nil, nil
}
