package node_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"testing"

	idtesting "github.com/agntcy/identity/internal/core/id/testing"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuertesting "github.com/agntcy/identity/internal/core/issuer/testing"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	coretesting "github.com/agntcy/identity/internal/core/testing"
	vctesting "github.com/agntcy/identity/internal/core/vc/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/stretchr/testify/assert"
)

func TestPublishVC(t *testing.T) {
	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	idRepo := idtesting.NewFakeIdRepository()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	vcRepo := vctesting.NewFakeVCRepository()
	sut := node.NewVerifiableCredentialService(verficationSrv, idRepo, issuerRepo, vcRepo)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	issuerRepo.CreateIssuer(context.Background(), issuer)
	credential := &vctypes.VerifiableCredential{
		ID: "VC_ID",
	}

	envelope, pubKey, err := signVCWithJose(credential)
	assert.NoError(t, err)

	resolverMD := &idtypes.ResolverMetadata{
		ID: fmt.Sprintf("DUO-%s", coretesting.ValidProofSub),
		VerificationMethod: []idtypes.VerificationMethod{
			{
				ID:           pubKey.KID,
				PublicKeyJwk: pubKey,
			},
		},
	}
	idRepo.CreateID(context.Background(), resolverMD)

	t.Run("should not return errors for JOSE VC", func(t *testing.T) {
		err = sut.Publish(context.Background(), envelope, &vctypes.Proof{})

		assert.NoError(t, err)
	})
}

func signVCWithJose(vc *vctypes.VerifiableCredential) (*vctypes.EnvelopedCredential, *idtypes.Jwk, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	payload, err := json.Marshal(vc)
	if err != nil {
		return nil, nil, err
	}

	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, "KEY-ID")
	signed, err := jws.Sign(payload, jws.WithKey(jwa.RS256(), pk, jws.WithProtectedHeaders(hdrs)))
	if err != nil {
		return nil, nil, err
	}

	pubkey, _ := jwk.PublicRawKeyOf(pk)
	key, _ := jwk.Import(pubkey)
	key.Set(jwk.AlgorithmKey, jwa.RS256())
	keyAsJson, _ := json.Marshal(key)

	var k idtypes.Jwk
	_ = json.Unmarshal(keyAsJson, &k)
	k.KID, _ = hdrs.KeyID()

	return &vctypes.EnvelopedCredential{
		EnvelopeType: vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE,
		Value:        string(signed),
	}, &k, nil
}
