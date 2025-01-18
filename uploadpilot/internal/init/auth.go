package init

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/quasoft/memstore"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/pkg/auth"
)

func initAuth(config *config.Config) error {
	store := memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	)
	gothic.Store = store

	goth.UseProviders(
		google.New(config.GoogleClientID, config.GoogleClientSecret, config.GoogleCallbackURL),
		github.New(config.GithubClientID, config.GithubClientSecret, config.GithubCallbackURL),
	)
	auth.SetSecretKey(config.JWTSecretKey)
	return nil
}
