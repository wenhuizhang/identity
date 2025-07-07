// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"github.com/agntcy/identity/pkg/jwk"
)

// VerificationMethod expresses verification methods, such as cryptographic
// public keys, which can be used to authenticate or authorize interactions
// with the entities represented by the ID. It is a part of the ResolverMetadata.
type VerificationMethod struct {
	// A unique id of the verification method.
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`

	// The public key used for the verification method.
	PublicKeyJwk *jwk.Jwk `json:"publicKeyJwk,omitempty" protobuf:"bytes,2,opt,name=public_key_jwk"`
}

// Service is used in ResolverMetadata to express ways of communicating with
// the node that published the document.
type Service struct {
	// ServiceEndpoint is a network address, such as an HTTP URL, of the
	// node.
	ServiceEndpoint []string `json:"serviceEndpoint" protobuf:"bytes,1,opt,name=service_endpoint"`
}

// ResolverMetadata represents a set of data describing the ID including mechanisms such as:
//   - cryptographic public keys - used to authenticate itself and prove
//     association with the ID
//   - service - ways of communicating with the node that published the document
//
// A ResolverMetadata can be retrieved by resolving an ID.
type ResolverMetadata struct {
	// The ID
	// The metadata below is related as claims to the ID
	ID string `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`

	// VerificationMethod is a list of cryptographic public keys, which can be used
	// to authenticate or authorize interactions with the entities represented by the ID.
	VerificationMethod []*VerificationMethod `json:"verificationMethod,omitempty" protobuf:"bytes,2,opt,name=verification_method"` //nolint:lll // Allow long lines

	// Service is used in ResolverMetadatas to express ways of communicating with
	// the node that published the document.
	Service []*Service `json:"service,omitempty" protobuf:"bytes,3,opt,name=service"`

	// AssertionMethod is used to specify how the entity represented by the ID
	// is expected to express claims, such as for the purposes of issuing a VCs.
	AssertionMethod []string `json:"assertionMethod,omitempty" protobuf:"bytes,4,opt,name=assertion_method"`

	// A controller is an entity that is authorized to make changes to a Resolver Metadata.
	Controller string
}

func (r *ResolverMetadata) GetJwks() *jwk.Jwks {
	jwks := jwk.Jwks{}

	for _, vm := range r.VerificationMethod {
		if vm.PublicKeyJwk != nil {
			jwks.Keys = append(jwks.Keys, vm.PublicKeyJwk)
		}
	}

	return &jwks
}
