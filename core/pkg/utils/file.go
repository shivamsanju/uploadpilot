package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

// ConvertToS3CompatibleFilename converts any filename to an S3-compatible filename
// respecting the 32 character limit while preserving the file extension.
// S3 allows alphanumeric characters, hyphens, underscores, and periods.
func ConvertToS3CompatibleFilename(filename string) string {
	// Extract file extension
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	// Remove invalid characters (keep only alphanumeric, hyphen, underscore, and period)
	reg := regexp.MustCompile("[^a-zA-Z0-9\\-_\\.]")
	sanitizedBaseName := reg.ReplaceAllString(baseName, "-")

	// Handle maximum length constraint (32 chars including the extension)
	maxBaseNameLength := 32 - len(ext)

	// If sanitized name is too long, simply truncate it
	if len(sanitizedBaseName) > maxBaseNameLength {
		sanitizedBaseName = sanitizedBaseName[:maxBaseNameLength]
	}

	return sanitizedBaseName + ext
}
