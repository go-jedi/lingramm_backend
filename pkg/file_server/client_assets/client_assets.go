package clientassets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	clientassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/apperrors"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
	"github.com/h2non/bimg"
)

// error definitions.
var (
	ErrFileTooLarge     = errors.New("file size exceeds maximum allowed limit")
	ErrInvalidPath      = errors.New("invalid file path")
	ErrDirectoryMissing = errors.New("directory does not exist")
)

// IClientAssets defines the interface for the file server.
//
//go:generate mockery --name=IClientAssets --output=mocks --case=underscore
type IClientAssets interface {
	UploadAndConvertToWebP(ctx context.Context, fileHeader *multipart.FileHeader) (clientassets.UploadAndConvertToWebpResponse, error)
}

type ClientAssets struct {
	url          string      // Url for connect to image from client
	dir          string      // Base directory for file storage
	maxFileSize  int64       // Maximum allowed file size
	imageQuality int         // Quality for image conversion
	dirPerm      os.FileMode // Permission mode for directories
	filePerm     os.FileMode // Permission mode for files

	uuid *uuid.UUID
}

// New creates a new ClientAssets instance with the given configuration.
func New(cfg config.FileServerConfig, uuid *uuid.UUID) *ClientAssets {
	fi := &ClientAssets{
		url:          cfg.ClientAssets.URL,
		dir:          cfg.ClientAssets.Dir,
		maxFileSize:  cfg.ClientAssets.MaxFileSize,
		imageQuality: cfg.ClientAssets.ImageQuality,
		dirPerm:      os.FileMode(cfg.DirPerm),
		filePerm:     os.FileMode(cfg.FilePerm),
		uuid:         uuid,
	}

	fi.init()

	return fi
}

// init sets default values for any unconfigured ClientAssets properties.
func (ca *ClientAssets) init() {
	if ca.maxFileSize == 0 {
		ca.maxFileSize = defaultMaxSize
	}

	if ca.imageQuality == 0 {
		ca.imageQuality = defaultQuality
	}

	if ca.dirPerm == 0 {
		ca.dirPerm = defaultDirPerm
	}

	if ca.filePerm == 0 {
		ca.filePerm = defaultFilePerm
	}
}

// getFileExt get file extension.
func (ca *ClientAssets) getFileExt(filename string) string {
	return strings.TrimSuffix(filepath.Ext(filename), filepath.Base(filename))
}

// sanitizePath checks if the given path is safe to use.
// returns cleaned path or error if path is invalid.
func (ca *ClientAssets) sanitizePath(path string) (string, error) {
	if path == "" {
		return "", ErrInvalidPath
	}

	const (
		parentDir       = ".."
		parentDirPrefix = ".." + string(filepath.Separator) // "../" or "..\"
		parentDirInPath = string(filepath.Separator) + ".." // "/.." or "\.."
	)

	// check for path traversal attempts.
	switch {
	case strings.HasPrefix(path, parentDirPrefix),
		strings.Contains(path, parentDirInPath),
		path == parentDir:
		return "", ErrInvalidPath
	}

	return filepath.Clean(path), nil
}

// validateFile checks if the file meets size and format requirements.
func (ca *ClientAssets) validateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > ca.maxFileSize {
		return ErrFileTooLarge
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if _, ok := clientassets.SupportedImageTypes[contentType]; !ok {
		log.Printf("unsupported file type: %s", contentType)
		return apperrors.ErrUnsupportedFormat
	}

	return nil
}

// ensureUploadDirectory checks if the target directory exists.
func (ca *ClientAssets) ensureUploadDirectory(uploadPath string) error {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrDirectoryMissing, uploadPath)
	}
	return nil
}

// readFileData reads file content with context support and size limitation.
func (ca *ClientAssets) readFileData(ctx context.Context, file multipart.File) ([]byte, error) {
	done := make(chan struct{})
	var data []byte
	var readErr error

	go func() {
		data, readErr = io.ReadAll(io.LimitReader(file, ca.maxFileSize))
		close(done)
	}()

	select {
	case <-done:
		return data, readErr
	case <-ctx.Done():
		return nil, ctx.Err() // return if context is canceled.
	}
}

// convertToWebP converts image data to WebP format with context support.
func (ca *ClientAssets) convertToWebP(ctx context.Context, rawFile []byte) ([]byte, error) {
	done := make(chan struct{})
	var webp []byte
	var convertErr error

	go func() {
		options := bimg.Options{
			Quality: ca.imageQuality,
			Type:    bimg.WEBP,
		}
		webp, convertErr = bimg.Resize(rawFile, options)
		close(done)
	}()

	select {
	case <-done:
		return webp, convertErr
	case <-ctx.Done():
		return nil, ctx.Err() // return if context is canceled.
	}
}

// UploadAndConvertToWebP handles the file upload process including validation, conversion and storage.
func (ca *ClientAssets) UploadAndConvertToWebP(ctx context.Context, fileHeader *multipart.FileHeader) (clientassets.UploadAndConvertToWebpResponse, error) {
	// check if context is already canceled.
	if err := ctx.Err(); err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, err
	}

	// measure upload duration.
	start := time.Now()
	defer func() {
		log.Printf("upload took %v", time.Since(start))
	}()

	// validate the uploaded file.
	if err := ca.validateFile(fileHeader); err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, err
	}

	// sanitize the target directory path.
	sanitizedDir, err := ca.sanitizePath(ca.dir)
	if err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("invalid directory: %w", err)
	}

	// open the uploaded file.
	file, err := fileHeader.Open()
	if err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	// prepare full upload path and verify directory exists.
	if err := ca.ensureUploadDirectory(sanitizedDir); err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, err
	}

	// read file data with context support.
	rawFile, err := ca.readFileData(ctx, file)
	if err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to read file: %w", err)
	}

	// convert image to WebP format.
	webp, err := ca.convertToWebP(ctx, rawFile)
	if err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to convert image: %w", err)
	}

	// generate unique filename and save the converted image.
	newName, err := ca.uuid.Generate()
	if err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to generate uuid: %w", err)
	}

	newFilePath := filepath.Join(sanitizedDir, newName+webpExt)

	if err := os.WriteFile(newFilePath, webp, ca.filePerm); err != nil {
		return clientassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to save converted image: %w", err)
	}

	return clientassets.UploadAndConvertToWebpResponse{
		NameFile:       newName + webpExt,
		ServerPathFile: filepath.Join(sanitizedDir, newName+webpExt),
		ClientPathFile: filepath.Join(ca.url, newName+webpExt),
		Extension:      webpExt,
		Quality:        ca.imageQuality,
		OldNameFile:    fileHeader.Filename,
		OldExtension:   ca.getFileExt(fileHeader.Filename),
	}, nil
}
