// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault/data"
)

type VaultService interface {
	ConnectVault() (*internalIssuerTypes.Vault, error)
	ListVaultIds() ([]string, error)
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

func (s *vaultService) ConnectVault() (*internalIssuerTypes.Vault, error) {

	vault, err := s.vaultRepository.ConnectVault()
	if err != nil {
		return nil, err
	}

	return vault, nil
}

func (s *vaultService) ListVaultIds() ([]string, error) {

	vaultIds, err := s.vaultRepository.ListVaultIds()
	if err != nil {
		return nil, err
	}

	return vaultIds, nil
}

func (s *vaultService) GetVault(vaultId string) (*internalIssuerTypes.Vault, error) {

	vault, err := s.vaultRepository.GetVault(vaultId)
	if err != nil {
		return nil, err
	}

	return vault, nil
}

func (s *vaultService) ForgetVault(vaultId string) error {

	err := s.vaultRepository.ForgetVault(vaultId)
	if err != nil {
		return err
	}

	return nil
}
