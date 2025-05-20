// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"github.com/agntcy/identity/internal/issuer/badge/data"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

type BadgeService interface {
	IssueBadge(vaultId, issuerId, metadataId, badgeValueFilePath string) (string, error)
	PublishBadge(
		vaultId, issuerId, metadataId string, badge *coreV1alpha.EnvelopedCredential,
	) (*coreV1alpha.EnvelopedCredential, error)
	ListBadgeIds(vaultId, issuerId, metadataId string) ([]string, error)
	GetBadge(vaultId, issuerId, metadataId, badgeId string) (*coreV1alpha.EnvelopedCredential, error)
	ForgetBadge(vaultId, issuerId, metadataId, badgeId string) error
}

type badgeService struct {
	badgeRepository data.BadgeRepository
}

func NewBadgeService(
	badgeRepository data.BadgeRepository,
) BadgeService {
	return &badgeService{
		badgeRepository: badgeRepository,
	}
}

func (s *badgeService) IssueBadge(vaultId, issuerId, metadataId, badgeValueFilePath string) (string, error) {
	badgeId, err := s.badgeRepository.IssueBadge(vaultId, issuerId, metadataId, badgeValueFilePath)
	if err != nil {
		return "", err
	}

	return badgeId, nil
}

func (s *badgeService) PublishBadge(
	vaultId, issuerId, metadataId string, badge *coreV1alpha.EnvelopedCredential,
) (*coreV1alpha.EnvelopedCredential, error) {
	badge, err := s.badgeRepository.PublishBadge(vaultId, issuerId, metadataId, badge)
	if err != nil {
		return nil, err
	}

	return badge, nil
}

func (s *badgeService) ListBadgeIds(vaultId, issuerId, metadataId string) ([]string, error) {
	badgeIds, err := s.badgeRepository.ListBadgeIds(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	return badgeIds, nil
}

func (s *badgeService) GetBadge(
	vaultId, issuerId, metadataId, badgeId string,
) (*coreV1alpha.EnvelopedCredential, error) {
	badge, err := s.badgeRepository.GetBadge(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	return badge, nil
}

func (s *badgeService) ForgetBadge(vaultId, issuerId, metadataId, badgeId string) error {
	err := s.badgeRepository.ForgetBadge(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return err
	}

	return nil
}
