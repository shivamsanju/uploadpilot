package kms

import "crypto/sha256"

var EncryptionKey []byte

func Init(encryptionKey string) error {
	key, err := GetValidKey(encryptionKey)
	if err != nil {
		return err
	}
	EncryptionKey = key
	return nil
}

// GetValidKey Create a SHA-256 hash of the encryption key to make sure it's 32 bytes
func GetValidKey(encryptionKey string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
