package rbac

import (
	"fmt"

	pgadapter "github.com/casbin/casbin-pg-adapter"
)

func NewPgAdapter(postgresURI, environnent string) (*pgadapter.Adapter, error) {
	a, err := pgadapter.NewAdapter(postgresURI, fmt.Sprintf("casbin_%s", environnent))
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}
	return a, nil
}
