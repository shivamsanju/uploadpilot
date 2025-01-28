package pdf

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type extractPDFContentTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewExtractPDFContentTask() tasks.Task {
	return &extractPDFContentTask{
		BaseTask:   tasks.NewBaseTask(),
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
	}
}

func (t *extractPDFContentTask) Do(ctx context.Context) error {
	t.leb.Publish(events.NewLogEvent(ctx, t.Data.WorkspaceID, t.Data.UploadID, "extracting pdf content", models.UploadLogLevelInfo))
	return nil
}

// var tempDirLock sync.Mutex // Ensures thread-safe access to the temporary folder.

// func ExtractPDFContent(ctx context.Context, upload *models.Upload) (*proc.Output, error) {
// 	tempDirLock.Lock()
// 	defer tempDirLock.Unlock()

// 	// Create a temporary directory.
// 	tempDir, err := os.MkdirTemp("", "pdf-extract-*")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
// 	}
// 	defer os.RemoveAll(tempDir) // Clean up temp directory before returning.

// 	// Download the PDF from S3 to the temporary directory.
// 	tempPDFPath := filepath.Join(tempDir, "input.pdf")
// 	file, err := os.Create(tempPDFPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create temporary file: %w", err)
// 	}
// 	defer file.Close()

// 	obj, err := infra.S3Client.GetObject(ctx, &s3.GetObjectInput{
// 		Bucket: &config.S3BucketName,
// 		Key:    &upload.StoredFileName,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get object from S3: %w", err)
// 	}
// 	defer obj.Body.Close()

// 	if _, err = io.Copy(file, obj.Body); err != nil {
// 		return nil, fmt.Errorf("failed to copy S3 object to temporary file: %w", err)
// 	}

// 	// Extract content from the PDF.
// 	outputDir := filepath.Join(tempDir, "extracted")
// 	if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
// 		return nil, fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	if err := api.ExtractPagesFile(tempPDFPath, outputDir, nil, nil); err != nil {
// 		return nil, fmt.Errorf("failed to extract content from PDF: %w", err)
// 	}

// 	// Zip the extracted content.
// 	zipFilePath := filepath.Join(tempDir, "extracted.zip")
// 	zipFile, err := os.Create(zipFilePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create zip file: %w", err)
// 	}
// 	defer zipFile.Close()

// 	if err := zipDirectory(outputDir, zipFile); err != nil {
// 		return nil, fmt.Errorf("failed to zip extracted content: %w", err)
// 	}

// 	// Upload the zip file back to S3.
// 	uploader := s3manager.NewUploader(infra.AWSConfig)
// 	zipFile.Seek(0, io.SeekStart) // Reset file pointer to the beginning.

// 	uploadOutput, err := uploader.Upload(ctx, &s3manager.UploadInput{
// 		Bucket: &config.S3BucketName,
// 		Key:    aws.String(fmt.Sprintf("extracted/%s.zip", upload.StoredFileName)),
// 		Body:   zipFile,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to upload zip file to S3: %w", err)
// 	}

// 	// Return the S3 URL of the uploaded zip file.
// 	return &proc.Output{
// 		URL: uploadOutput.Location,
// 	}, nil
// }

// // zipDirectory compresses the contents of a directory into a zip file.
// func zipDirectory(sourceDir string, writer io.Writer) error {
// 	zipWriter := zip.NewWriter(writer)
// 	defer zipWriter.Close()

// 	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		relPath, err := filepath.Rel(sourceDir, path)
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() {
// 			return nil
// 		}

// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		header, err := zip.FileInfoHeader(info)
// 		if err != nil {
// 			return err
// 		}
// 		header.Name = relPath
// 		header.Method = zip.Deflate

// 		zipFile, err := zipWriter.CreateHeader(header)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = io.Copy(zipFile, file)
// 		return err
// 	})
// }
