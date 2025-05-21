// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault/data"
	"github.com/google/uuid"
)

type VaultService interface {
	ConnectVault(
		vaultType internalIssuerTypes.VaultType, config internalIssuerTypes.VaultConfig,
	) (*internalIssuerTypes.Vault, error)
	GetAllVaults() ([]*internalIssuerTypes.Vault, error)
	GetVault(vaultId string) (*internalIssuerTypes.Vault, error)
	ForgetVault(vaultId string) error
}

type vaultService struct {
	vaultRepository data.VaultRepository
}

func NewVaultService(
	vaultRepository data.VaultRepository,
) VaultService {
	return &vaultService{
		vaultRepository: vaultRepository,
	}
}

func (s *vaultService) ConnectVault(
	vaultType internalIssuerTypes.VaultType, config internalIssuerTypes.VaultConfig,
) (*internalIssuerTypes.Vault, error) {
	vault := internalIssuerTypes.Vault{
		Id:     uuid.NewString(),
		Type:   internalIssuerTypes.VaultTypeTxt,
		Config: config,
	}

	_, err := s.vaultRepository.AddVault(&vault)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func (s *vaultService) GetAllVaults() ([]*internalIssuerTypes.Vault, error) {
	vaults, err := s.vaultRepository.GetAllVaults()
	if err != nil {
		return nil, err
	}

	return vaults, nil
}

func (s *vaultService) GetVault(vaultId string) (*internalIssuerTypes.Vault, error) {
	vault, err := s.vaultRepository.GetVault(vaultId)
	if err != nil {
		return nil, err
	}

	return vault, nil
}

func (s *vaultService) ForgetVault(vaultId string) error {
	err := s.vaultRepository.RemoveVault(vaultId)
	if err != nil {
		return err
	}

	return nil
}
