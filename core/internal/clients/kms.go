package clients

import "github.com/uploadpilot/core/pkg/vault"

type KMSOpts struct {
	EncryptionKey string
}

func NewKMSClient(opts *KMSOpts) (vault.KMS, error) {
	return vault.NewKMS(opts.EncryptionKey)
}
