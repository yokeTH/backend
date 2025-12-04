package jwk

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type LocalKeyEncryptor struct {
	kek []byte
}

func NewLocalKeyEncryptorFromHex(keyHex string) (*LocalKeyEncryptor, error) {
	kek, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("decode KEK hex: %w", err)
	}

	switch len(kek) {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("invalid KEK length %d (must be 16, 24, or 32 bytes)", len(kek))
	}

	return &LocalKeyEncryptor{
		kek: kek,
	}, nil
}

func (e *LocalKeyEncryptor) Encrypt(ctx context.Context, kekRef string, plaintext []byte) (ciphertext, nonce, wrappedDEK []byte, err error) {
	_ = ctx

	block, err := aes.NewCipher(e.kek)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("new cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("new GCM: %w", err)
	}

	nonce = make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext = aead.Seal(nil, nonce, plaintext, nil)

	wrappedDEK = nil

	return ciphertext, nonce, wrappedDEK, nil
}

func (e *LocalKeyEncryptor) Decrypt(ctx context.Context, kekRef string, ciphertext, nonce []byte) ([]byte, error) {
	_ = ctx
	_ = kekRef

	block, err := aes.NewCipher(e.kek)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size: got %d, want %d", len(nonce), aead.NonceSize())
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
