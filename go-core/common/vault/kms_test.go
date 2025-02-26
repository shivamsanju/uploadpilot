package vault_test

import (
	"testing"

	"github.com/uploadpilot/go-core/common/vault"
)

func TestKMS(t *testing.T) {
	secret := "thisisaverysecureandlongsecretkey!"
	salt := "thisisaverysecureandsalt!"
	kms, err := vault.NewKMS(secret)
	if err != nil {
		t.Fatalf("failed to create KMS: %v", err)
	}

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		tests := []struct {
			name      string
			plaintext string
		}{
			{"Simple Text", "Hello, World!"},
			{"Empty String", ""},
			{"Special Characters", "!@#$%^&*()_+"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				encrypted, salt, err := kms.Encrypt(tt.plaintext)
				if err != nil {
					t.Fatalf("Encrypt failed: %v", err)
				}

				decrypted, err := kms.Decrypt(encrypted, salt)
				if err != nil {
					t.Fatalf("Decrypt failed: %v", err)
				}

				if decrypted != tt.plaintext {
					t.Errorf("expected %q, got %q", tt.plaintext, decrypted)
				}
			})
		}
	})

	t.Run("Hash and VerifyHash", func(t *testing.T) {
		tests := []struct {
			name      string
			plaintext string
		}{
			{"Simple Password", "password123"},
			{"Empty String", ""},
			{"Complex Password", "P@ssw0rd!2024"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hash, err := kms.Hash(tt.plaintext, salt)
				if err != nil {
					t.Fatalf("Hash failed: %v", err)
				}

				verified, err := kms.VerifyHash(tt.plaintext, hash, salt)
				if err != nil {
					t.Fatalf("VerifyHash failed: %v", err)
				}

				if !verified {
					t.Errorf("expected hash verification to succeed for %q", tt.plaintext)
				}
			})
		}
	})
}
