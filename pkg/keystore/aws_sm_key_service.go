// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keystore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/agntcy/identity/internal/pkg/ptrutil"
	"github.com/agntcy/identity/pkg/jwk"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smtypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

var errNilAwsConfig = errors.New("please provide a configuration for the AWS SM client")

type AwsSmKeyService struct {
	client *secretsmanager.Client
	cfg    *AwsSmStorageConfig
}

type AwsSmStorageConfig struct {
	AwsCfg      *aws.Config
	MountPath   string
	KeyBasePath string
	KmsKeyID    *string
}

type awsSMData struct {
	JWK *jwk.Jwk
}

func NewAwsSmKeyService(cfg *AwsSmStorageConfig) (KeyService, error) {
	if cfg == nil || cfg.AwsCfg == nil {
		return nil, errNilAwsConfig
	}

	return &AwsSmKeyService{
		client: secretsmanager.NewFromConfig(*cfg.AwsCfg),
		cfg:    cfg,
	}, nil
}

func (s *AwsSmKeyService) SaveKey(ctx context.Context, id string, priv *jwk.Jwk) error {
	if priv == nil {
		return errors.New("jwk cannot be nil")
	}

	data := &awsSMData{JWK: priv}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JWK: %w", err)
	}

	req := &secretsmanager.CreateSecretInput{
		Name:         ptrutil.Ptr(s.buildKeyPath(id)),
		KmsKeyId:     s.cfg.KmsKeyID,
		SecretBinary: bytes,
	}

	_, err = s.client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to write JWK to AWS Secrets Manager: %w", err)
	}

	return nil
}

func (s *AwsSmKeyService) RetrievePubKey(ctx context.Context, id string) (*jwk.Jwk, error) {
	priv, err := s.retrieveKey(ctx, id)
	if err != nil {
		return nil, err
	}

	return priv.PublicKey(), nil
}

func (s *AwsSmKeyService) RetrievePrivKey(ctx context.Context, id string) (*jwk.Jwk, error) {
	return s.retrieveKey(ctx, id)
}

func (s *AwsSmKeyService) DeleteKey(ctx context.Context, id string) error {
	_, err := s.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId: ptrutil.Ptr(s.buildKeyPath(id)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete key from Vault: %w", err)
	}

	return nil
}

func (s *AwsSmKeyService) ListKeys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)

	paginator := secretsmanager.NewListSecretsPaginator(s.client, &secretsmanager.ListSecretsInput{
		Filters: []smtypes.Filter{
			{
				Key:    "name",
				Values: []string{s.buildKeyPath("")},
			},
		},
	})
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list keys in Vault: %w", err)
		}

		for idx := range len(out.SecretList) {
			secret := &out.SecretList[idx]

			if secret.Name != nil {
				keys = append(keys, *secret.Name)
			}
		}
	}

	return keys, nil
}

func (s *AwsSmKeyService) retrieveKey(ctx context.Context, id string) (*jwk.Jwk, error) {
	secret, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: ptrutil.Ptr(s.buildKeyPath(id)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read JWK from AWS Secrets Manager: %w", err)
	}

	if secret == nil || len(secret.SecretBinary) == 0 {
		return nil, errors.New("key not found in AWS Secrets Manager")
	}

	var data awsSMData
	if err := json.Unmarshal(secret.SecretBinary, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWK: %w", err)
	}

	return data.JWK, nil
}

func (s *AwsSmKeyService) buildKeyPath(id string) string {
	return path.Join(s.cfg.MountPath, s.cfg.KeyBasePath, id)
}
