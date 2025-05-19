// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"
	"errors"

	"github.com/agntcy/identity/internal/core"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
)

const (
	ValidProofIssuer string = "ISSUER"
	ValidProofSub    string = "SUBJECT"
)

type FakeTruthyVerificationService struct{}

func NewFakeTruthyVerificationService() core.VerificationService {
	return &FakeTruthyVerificationService{}
}

func (f *FakeTruthyVerificationService) VerifyCommonName(
	ctx context.Context,
	commonName *string,
	proof *vctypes.Proof,
) error {
	return nil
}

func (f *FakeTruthyVerificationService) VerifyProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, string, error) {
	return ValidProofIssuer, ValidProofSub, nil
}

type FalsyProofVerificationServiceStub struct{}

func NewFalsyProofVerificationServiceStub() core.VerificationService {
	return &FalsyProofVerificationServiceStub{}
}

func (f *FalsyProofVerificationServiceStub) VerifyCommonName(
	ctx context.Context,
	commonName *string,
	proof *vctypes.Proof,
) error {
	panic("unimplemented")
}

func (f *FalsyProofVerificationServiceStub) VerifyProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, string, error) {
	return "", "", errors.New("UNPROVEN")
}

type FalsyCommonNameVerificationServiceStub struct{}

func NewFalsyCommonNameVerificationServiceStub() core.VerificationService {
	return &FalsyCommonNameVerificationServiceStub{}
}

func (f *FalsyCommonNameVerificationServiceStub) VerifyCommonName(
	ctx context.Context,
	commonName *string,
	proof *vctypes.Proof,
) error {
	return errors.New("UNPROVEN")
}

func (f *FalsyCommonNameVerificationServiceStub) VerifyProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, string, error) {
	return "", "", nil
}
