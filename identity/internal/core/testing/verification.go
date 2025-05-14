package testing

import (
	"context"

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
