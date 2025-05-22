// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"context"
	"errors"
	"net/url"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/converters"
	issuersdk "github.com/agntcy/identity/sdk/node-go/client/issuer_service"
	"github.com/agntcy/identity/sdk/node-go/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type NodeClientProvider interface {
	New(host string) (NodeClient, error)
}

type nodeClientProvider struct{}

func NewNodeClientProvider() NodeClientProvider {
	return &nodeClientProvider{}
}

func (nodeClientProvider) New(host string) (NodeClient, error) {
	return NewNodeClient(host)
}

type NodeClient interface {
	RegisterIssuer(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) error
}

type nodeClient struct {
	issuer issuersdk.ClientService
}

func NewNodeClient(host string) (NodeClient, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &nodeClient{
		issuer: issuersdk.New(httptransport.New(u.Host, "", []string{u.Scheme}), strfmt.Default),
	}, nil
}

func (c *nodeClient) RegisterIssuer(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) error {
	resp, err := c.issuer.RegisterIssuer(&issuersdk.RegisterIssuerParams{
		Body: &models.V1alpha1RegisterIssuerRequest{
			Issuer: &models.V1alpha1Issuer{
				CommonName:      issuer.CommonName,
				Organization:    issuer.Organization,
				SubOrganization: issuer.SubOrganization,
				PublicKey:       converters.Convert[models.V1alpha1Jwk](issuer.PublicKey.PublicKey()),
			},
			Proof: &models.V1alpha1Proof{
				Type:       proof.Type,
				ProofValue: proof.ProofValue,
			},
		},
	})
	if err != nil {
		return err
	}

	if resp.Code() > 299 {
		return errors.New("unsuccessful")
	}

	return nil
}
