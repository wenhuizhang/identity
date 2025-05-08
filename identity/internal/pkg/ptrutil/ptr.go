// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package ptrutil

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}
