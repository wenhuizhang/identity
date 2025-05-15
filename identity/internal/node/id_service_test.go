// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node_test

import (
	"context"
	"testing"

	idtesting "github.com/agntcy/identity/internal/core/id/testing"
	issuertesting "github.com/agntcy/identity/internal/core/issuer/testing"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	coretesting "github.com/agntcy/identity/internal/core/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	idRepo := idtesting.NewFakeIdRepository()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIdService(verficationSrv, idRepo, issuerRepo)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	issuerRepo.CreateIssuer(context.Background(), issuer)

	t.Run("should not return an error", func(t *testing.T) {
		_, err := sut.Generate(t.Context(), issuer, &vctypes.Proof{})

		assert.NoError(t, err)
	})
}
