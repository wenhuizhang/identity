// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"errors"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
)

func VerifyCredential(credential *vctypes.VerifiableCredential) (bool, error) {

	return false, errors.New("VerifyCredential not implemented yet")
}
