// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"fmt"

	nodeapi "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/agntcy/identity/internal/pkg/converters"
	"github.com/agntcy/identity/internal/pkg/grpcutil"
	"google.golang.org/protobuf/types/known/emptypb"
)

type vcService struct {
	vcSrv node.VerifiableCredentialService
}

func NewVcService(vcSrv node.VerifiableCredentialService) nodeapi.VcServiceServer {
	return &vcService{
		vcSrv: vcSrv,
	}
}

// Publish an issued Verifiable Credential
func (s *vcService) Publish(
	ctx context.Context,
	req *nodeapi.PublishRequest,
) (*emptypb.Empty, error) {
	err := s.vcSrv.Publish(
		ctx,
		converters.Convert[vctypes.EnvelopedCredential](req.Vc),
		converters.Convert[vctypes.Proof](req.Proof),
	)
	if err != nil {
		if errtypes.IsErrorInfo(err, errtypes.ERROR_REASON_INTERNAL) {
			return nil, grpcutil.InternalError(err)
		}

		return nil, grpcutil.BadRequestError(err)
	}

	return &emptypb.Empty{}, nil
}

// Verify an existing Verifiable Credential
func (vcService) Verify(
	ctx context.Context,
	req *nodeapi.VerifyRequest,
) (*emptypb.Empty, error) {
	return nil, fmt.Errorf("not implemented")
}

// Returns the well-known Verifiable Credentials for the specified Id
func (vcService) GetWellKnown(
	ctx context.Context,
	req *nodeapi.GetVcWellKnownRequest,
) (*nodeapi.GetVcWellKnownResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// Search for Verifiable Credentials based on the specified criteria
func (vcService) Search(
	ctx context.Context,
	req *nodeapi.SearchRequest,
) (*nodeapi.SearchResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
