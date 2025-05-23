// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Cache struct {
	VaultId    string `json:"vaultId,omitempty"`
	IssuerId   string `json:"issuerId,omitempty"`
	MetadataId string `json:"metadata,omitempty"`
	BadgeId    string `json:"badgeId,omitempty"`
	KeyID      string `json:"kid,omitempty"`
}

// getCacheFile returns the path to the cache file
func getCacheFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".identity", "cache.json"), nil
}

// SaveCache saves the cache to a file
func SaveCache(cache *Cache) error {
	cacheFile, err := getCacheFile()
	if err != nil {
		return err
	}

	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cache); err != nil {
		return err
	}

	return nil
}

// LoadCache loads the cache from a file, creating the file if it doesn't exist and returning an empty cache
func LoadCache() (*Cache, error) {
	cacheFile, err := getCacheFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return &Cache{}, nil
	}

	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache Cache

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// ClearCache clears the cache file
func ClearCache() error {
	cacheFile, err := getCacheFile()
	if err != nil {
		return err
	}

	if err := os.Remove(cacheFile); err != nil {
		return err
	}

	return nil
}
