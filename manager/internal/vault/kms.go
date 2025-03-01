package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

type KMS interface {
	Encrypt(plaintext string) (value string, salt string, err error)
	Decrypt(value string, salt string) (plaintext string, err error)
	Hash(plaintext string, salt string) (value string, err error)
	VerifyHash(plaintext, value string, salt string) (bool, error)
}

type kms struct {
	secret []byte
	salt   []byte
}

func NewKMS(secret string) (KMS, error) {
	if len(secret) < 32 {
		return nil, errors.New("secret key must be at least 32 bytes long")
	}

	return &kms{
		secret: []byte(secret)[:32],
	}, nil
}

func (k *kms) Encrypt(plaintext string) (string, string, error) {
	block, err := aes.NewCipher(k.secret)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(nonce), nil
}

func (k *kms) Decrypt(value string, salt string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(k.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (k *kms) Hash(plaintext string, salt string) (string, error) {
	hashedValue := k.hash(plaintext, []byte(salt))
	return base64.StdEncoding.EncodeToString(hashedValue), nil
}

func (k *kms) hash(plaintext string, salt []byte) []byte {
	return argon2.IDKey([]byte(plaintext), salt, 1, 64*1024, 4, 32)
}

func (k *kms) VerifyHash(plaintext, value string, salt string) (bool, error) {
	hashedValue := k.hash(plaintext, []byte(salt))
	return base64.StdEncoding.EncodeToString(hashedValue) == value, nil
}
