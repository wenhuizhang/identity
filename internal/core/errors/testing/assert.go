// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"testing"

	"github.com/agntcy/identity/internal/core/errors/types"
	"github.com/stretchr/testify/assert"
)

func AssertErrorInfoReason(t *testing.T, err error, reason types.ErrorReason) {
	t.Helper()

	var errInfo types.ErrorInfo
	assert.ErrorAs(t, err, &errInfo)
	assert.Equal(t, reason, errInfo.Reason)
}
