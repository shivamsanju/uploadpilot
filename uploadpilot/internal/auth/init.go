package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/quasoft/memstore"
	"github.com/uploadpilot/uploadpilot/internal/config"
)

var (
	secretKey []byte
)

func Init() error {
	store := memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	)
	gothic.Store = store

	goth.UseProviders(
		google.New(config.GoogleClientID, config.GoogleClientSecret, config.GoogleCallbackURL, "email", "profile"),
		github.New(config.GithubClientID, config.GithubClientSecret, config.GithubCallbackURL, "email", "user"),
	)
	secretKey = []byte(config.JWTSecretKey)
	return nil
}
