// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keystore

import (
	"errors"
	"fmt"
)

type StorageType int

const (
	FileStorage StorageType = iota
	OnePasswordStorage
)

func (s StorageType) String() string {
	return [...]string{"file", "1password"}[s]
}

type FileStorageConfig struct {
	FilePath string
}

func NewKeyService(storageType StorageType, config interface{}) (KeyService, error) {
	switch storageType {
	case FileStorage:
		c, err := getConfig[FileStorageConfig](config)
		if err != nil {
			return nil, err
		}

		return &LocalFileKeyService{FilePath: c.FilePath}, nil

	case OnePasswordStorage:
		// Commented out until implemented
		// c, err := getConfig[OnePasswordConfig](config)
		// if err != nil {
		//     return nil, err
		// }
		// return &OnePasswordKeyService{...}, nil
		return nil, errors.New("1password storage not implemented")
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

func getConfig[T any](config interface{}) (T, error) {
	var zero T

	if config == nil {
		return zero, errors.New("nil config provided")
	}

	if c, ok := config.(T); ok {
		return c, nil
	}

	return zero, fmt.Errorf("invalid config type: expected %T", zero)
}
