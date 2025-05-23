// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package nodeapi

import (
	"context"
	"net/url"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/converters"
	idsdk "github.com/agntcy/identity/sdk/node-go/client/id_service"
	issuersdk "github.com/agntcy/identity/sdk/node-go/client/issuer_service"
	vcsdk "github.com/agntcy/identity/sdk/node-go/client/vc_service"
	"github.com/agntcy/identity/sdk/node-go/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type ClientProvider interface {
	New(host string) (NodeClient, error)
}

type clientProvider struct{}

func NewNodeClientProvider() ClientProvider {
	return &clientProvider{}
}

func (clientProvider) New(host string) (NodeClient, error) {
	return NewNodeClient(host)
}

type NodeClient interface {
	RegisterIssuer(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) error
}

type nodeClient struct {
	id     idsdk.ClientService
	issuer issuersdk.ClientService
	vc     vcsdk.ClientService
}

func NewNodeClient(host string) (NodeClient, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &nodeClient{
		id:     idsdk.New(httptransport.New(u.Host, "", []string{u.Scheme}), strfmt.Default),
		issuer: issuersdk.New(httptransport.New(u.Host, "", []string{u.Scheme}), strfmt.Default),
		vc:     vcsdk.New(httptransport.New(u.Host, "", []string{u.Scheme}), strfmt.Default),
	}, nil
}

func (c *nodeClient) RegisterIssuer(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) error {
	_, err := c.issuer.RegisterIssuer(&issuersdk.RegisterIssuerParams{
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

	return nil
}
