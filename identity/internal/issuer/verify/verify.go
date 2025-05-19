// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"errors"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

func VerifyCredential(credential *coreV1alpha.VerifiableCredential) (bool, error) {
	return false, errors.New("VerifyCredential not implemented yet")
}
