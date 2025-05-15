// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"errors"
	"fmt"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	issuerSetup "github.com/agntcy/identity/internal/issuer/setup"
)

func VerifyCredential(credential *vctypes.VerifiableCredential) (bool, error) {
	_, err := issuerSetup.ReadNetworkConfig()
	if err != nil {
		return false, fmt.Errorf("error reading network config: %w", err)
	}

	return false, errors.New("VerifyCredential not implemented yet")
}
