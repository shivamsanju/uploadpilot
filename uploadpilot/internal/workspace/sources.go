package workspace

import "github.com/uploadpilot/uploadpilot/internal/db/models"

var AllowedSources = []models.AllowedSources{
	models.FileUpload,
	models.Webcamera,
	models.Audio,
	models.ScreenCapture,
	models.Box,
	models.Dropbox,
	models.Facebook,
	models.GoogleDrive,
	models.GooglePhotos,
	models.Instagram,
	models.OneDrive,
	models.Unsplash,
	models.Url,
	models.Zoom,
}
