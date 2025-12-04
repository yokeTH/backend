package jwk

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type KeyEncryptor interface {
	Encrypt(ctx context.Context, kekRef string, plaintext []byte) (ciphertext, nonce, wrappedDEK []byte, err error)
	Decrypt(ctx context.Context, kekRef string, ciphertext, nonce []byte) ([]byte, error)
}

type JWKGenerator struct {
	encryptor    KeyEncryptor
	kekRef       string
	alg          string
	defaultState model.JWKKeyStatus
	ttl          time.Duration
	now          func() time.Time
}

func NewJWKGenerator(encryptor KeyEncryptor, kekRef string, defaultState model.JWKKeyStatus, ttl time.Duration) *JWKGenerator {
	return &JWKGenerator{
		encryptor:    encryptor,
		kekRef:       kekRef,
		alg:          "ES256",
		defaultState: defaultState,
		ttl:          ttl,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

type ecThumbprintJWK struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type ecPublicJWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Alg string `json:"alg,omitempty"`
	Use string `json:"use,omitempty"`
	Kid string `json:"kid,omitempty"`
}

type ecPrivateJWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	D   string `json:"d"`
	Alg string `json:"alg,omitempty"`
	Use string `json:"use,omitempty"`
	Kid string `json:"kid,omitempty"`
}

func encodeBigIntToBase64URL(v *big.Int, size int) string {
	b := v.Bytes()
	if len(b) < size {
		padded := make([]byte, size)
		copy(padded[size-len(b):], b)
		b = padded
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func computeECThumbprintRFC7638(x, y string) (string, error) {
	t := ecThumbprintJWK{
		Crv: "P-256",
		Kty: "EC",
		X:   x,
		Y:   y,
	}
	j, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("marshal thumbprint JWK: %w", err)
	}

	sum := sha256.Sum256(j)
	return base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func (g *JWKGenerator) Generate(ctx context.Context) (*model.JWKKeyModel, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate EC key: %w", err)
	}

	curveParams := privKey.Curve.Params()
	size := (curveParams.BitSize + 7) / 8 // 32 bytes for P-256

	xB64 := encodeBigIntToBase64URL(privKey.X, size)
	yB64 := encodeBigIntToBase64URL(privKey.Y, size)
	dB64 := encodeBigIntToBase64URL(privKey.D, size)

	kid, err := computeECThumbprintRFC7638(xB64, yB64)
	if err != nil {
		return nil, fmt.Errorf("compute thumbprint: %w", err)
	}

	pubJWK := ecPublicJWK{
		Kty: "EC",
		Crv: "P-256",
		X:   xB64,
		Y:   yB64,
		Alg: g.alg,
		Use: "sig",
		Kid: kid, //RFC7638
	}

	privJWK := ecPrivateJWK{
		Kty: "EC",
		Crv: "P-256",
		X:   xB64,
		Y:   yB64,
		D:   dB64,
		Alg: g.alg,
		Use: "sig",
		Kid: kid,
	}

	pubJSON, err := json.Marshal(pubJWK)
	if err != nil {
		return nil, fmt.Errorf("marshal public JWK: %w", err)
	}

	privJSON, err := json.Marshal(privJWK)
	if err != nil {
		return nil, fmt.Errorf("marshal private JWK: %w", err)
	}

	privCiphertext, privNonce, wrappedDEK, err := g.encryptor.Encrypt(ctx, g.kekRef, privJSON)
	if err != nil {
		return nil, fmt.Errorf("encrypt private JWK: %w", err)
	}

	now := g.now()

	var notBefore *time.Time
	var notAfter *time.Time
	var notAfterEpoch *int64

	nb := now
	notBefore = &nb

	if g.ttl > 0 {
		na := now.Add(g.ttl).UTC()
		notAfter = &na
		epoch := na.Unix()
		notAfterEpoch = &epoch
	}

	m := &model.JWKKeyModel{
		KID:    kid,
		Alg:    g.alg,
		Status: g.defaultState,

		PublicJWK: string(pubJSON),

		PrivCiphertext: privCiphertext,
		PrivNonce:      privNonce,
		WrappedDEK:     wrappedDEK,

		KEKRef: g.kekRef,

		CreatedAt: now,
		RotatedAt: nil,
		NotBefore: notBefore,
		NotAfter:  notAfter,

		NotAfterEpoch: notAfterEpoch,
	}

	return m, nil
}
