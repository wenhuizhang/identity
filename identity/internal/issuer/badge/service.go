// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"log"

	"github.com/agntcy/identity/internal/issuer/badge/data"
	"github.com/google/uuid"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type BadgeService interface {
	IssueBadge(vaultId, issuerId, metadataId, badgeContent string) (string, error)
	PublishBadge(
		vaultId, issuerId, metadataId string, badge *internalIssuerTypes.Badge,
	) (*internalIssuerTypes.Badge, error)
	GetAllBadges(vaultId, issuerId, metadataId string) ([]*internalIssuerTypes.Badge, error)
	GetBadge(vaultId, issuerId, metadataId, badgeId string) (*internalIssuerTypes.Badge, error)
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

func (s *badgeService) IssueBadge(vaultId, issuerId, metadataId, badgeContent string) (string, error) {

	envelopedCredential := coreV1alpha.EnvelopedCredential{
		EnvelopeType: coreV1alpha.CredentialEnvelopeType_CREDENTIAL_ENVELOPE_TYPE_JOSE.Enum(),
		Value:        &badgeContent,
	}

	badge := internalIssuerTypes.Badge{
		Id:                  uuid.New().String(),
		EnvelopedCredential: &envelopedCredential,
	}

	badgeId, err := s.badgeRepository.AddBadge(vaultId, issuerId, metadataId, &badge)
	if err != nil {
		return "", err
	}

	return badgeId, nil
}

func (s *badgeService) PublishBadge(
	vaultId, issuerId, metadataId string, badge *internalIssuerTypes.Badge,
) (*internalIssuerTypes.Badge, error) {
	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	publishRequest := nodeV1alpha.PublishRequest{
		Vc:    badge.EnvelopedCredential,
		Proof: &proof,
	}

	log.Default().Println("Publishing badge with request: ", &publishRequest)

	return badge, nil
}

func (s *badgeService) GetAllBadges(vaultId, issuerId, metadataId string) ([]*internalIssuerTypes.Badge, error) {
	badges, err := s.badgeRepository.GetAllBadges(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	return badges, nil
}

func (s *badgeService) GetBadge(
	vaultId, issuerId, metadataId, badgeId string,
) (*internalIssuerTypes.Badge, error) {
	badge, err := s.badgeRepository.GetBadge(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	return badge, nil
}

func (s *badgeService) ForgetBadge(vaultId, issuerId, metadataId, badgeId string) error {
	err := s.badgeRepository.RemoveBadge(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return err
	}

	return nil
}
