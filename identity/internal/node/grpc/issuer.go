package grpc

import (
	"context"
	"fmt"

	nodeapi "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
)

type issuerService struct{}

func NewIssuerService() nodeapi.IssuerServiceServer {
	return &issuerService{}
}

// Register an issuer by providing the issuer details
func (issuerService) Register(
	ctx context.Context,
	req *nodeapi.RegisterIssuerRequest,
) (*nodeapi.RegisterIssuerResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// Returns the well-known document content for an issuer in
// Json Web Key Set (JWKS) format
func (issuerService) GetWellKnown(
	ctx context.Context,
	req *nodeapi.GetIssuerWellKnownRequest,
) (*nodeapi.GetIssuerWellKnownResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
