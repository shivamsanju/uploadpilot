package driver

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractDSNFromURI(dbURI string) (string, error) {
	parsedURL, err := url.Parse(dbURI)
	if err != nil {
		return "", fmt.Errorf("invalid database URI: %w", err)
	}

	switch parsedURL.Scheme {
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			parsedURL.Hostname(),
			parsedURL.User.Username(),
			getPassword(parsedURL),
			strings.TrimPrefix(parsedURL.Path, "/"),
			parsedURL.Port(),
			parsedURL.Query().Get("sslmode"),
		), nil

	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
			parsedURL.User.Username(),
			getPassword(parsedURL),
			parsedURL.Host,
			strings.TrimPrefix(parsedURL.Path, "/"),
			parsedURL.RawQuery,
		), nil

	case "sqlite", "sqlite3":
		return parsedURL.Path, nil

	default:
		return "", fmt.Errorf("unsupported database type: %s", parsedURL.Scheme)
	}
}

// getPassword extracts the password from the URL.User field
func getPassword(parsedURL *url.URL) string {
	if pass, hasPass := parsedURL.User.Password(); hasPass {
		return pass
	}
	return ""
}
