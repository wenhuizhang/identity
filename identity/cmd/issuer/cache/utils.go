// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package cache

import "fmt"

//nolint:lll // Allow long lines for CLI
func (c *Cache) ValidateVaultId() error {
	if c == nil || c.VaultId == "" {
		return fmt.Errorf(
			"no vault found in the local configuration. Please load an existing vault or connect to a new vault",
		)
	}

	return nil
}

//nolint:lll // Allow long lines for CLI
func (c *Cache) ValidateKeyId() error {
	if c.KeyID == "" {
		return fmt.Errorf(
			"no key found in the local configuration. Please load an existing key from the vault or generate a new key",
		)
	}

	return nil
}

//nolint:lll // Allow long lines for CLI
func (c *Cache) ValidateIssuerId() error {
	if c.IssuerId == "" {
		return fmt.Errorf(
			"no issuer found in the local configuration. Please load an existing issuer or register a new issuer",
		)
	}

	return nil
}

//nolint:lll // Allow long lines for CLI
func (c *Cache) ValidateMetadataId() error {
	if c.MetadataId == "" {
		return fmt.Errorf(
			"no metadata found in the local configuration. Please load an existing metadata or generate a new metadata",
		)
	}

	return nil
}

func (c *Cache) ValidateBadgeId() error {
	if c.BadgeId == "" {
		return fmt.Errorf(
			"no badge found in the local configuration. Please load an existing badge or issue a new badge",
		)
	}

	return nil
}

func (c *Cache) Validate() error {
	if err := c.ValidateForBadge(); err != nil {
		return err
	}

	if err := c.ValidateBadgeId(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) ValidateForBadge() error {
	if err := c.ValidateForMetadata(); err != nil {
		return err
	}

	if err := c.ValidateMetadataId(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) ValidateForMetadata() error {
	if err := c.ValidateForIssuer(); err != nil {
		return err
	}

	if err := c.ValidateIssuerId(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) ValidateForIssuer() error {
	if err := c.ValidateForKey(); err != nil {
		return err
	}

	if err := c.ValidateKeyId(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) ValidateForKey() error {
	if err := c.ValidateVaultId(); err != nil {
		return err
	}

	return nil
}
