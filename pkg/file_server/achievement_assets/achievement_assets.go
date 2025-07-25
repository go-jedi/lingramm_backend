package achievementassets

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
	achievementassets "github.com/go-jedi/lingramm_backend/internal/domain/file_server/achievement_assets"
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

// IAchievementAssets defines the interface for the file server.
//
//go:generate mockery --name=IAchievementAssets --output=mocks --case=underscore
type IAchievementAssets interface {
	UploadAndConvertToWebP(ctx context.Context, fileHeader *multipart.FileHeader) (achievementassets.UploadAndConvertToWebpResponse, error)
}

type AchievementAssets struct {
	maxFileSize  int64 // Maximum allowed file size
	uuid         *uuid.UUID
	url          string      // URL for client access
	dir          string      // File storage directory
	imageQuality int         // Image quality after conversion
	dirPerm      os.FileMode // Directory permission mode (uint32)
	filePerm     os.FileMode // File permission mode (uint32)
}

// New creates a new AchievementAssets instance with the given configuration.
func New(cfg config.FileServerConfig, uuid *uuid.UUID) *AchievementAssets {
	aa := &AchievementAssets{
		url:          cfg.AchievementAssets.URL,
		dir:          cfg.AchievementAssets.Dir,
		maxFileSize:  cfg.AchievementAssets.MaxFileSize,
		imageQuality: cfg.AchievementAssets.ImageQuality,
		dirPerm:      os.FileMode(cfg.DirPerm),
		filePerm:     os.FileMode(cfg.FilePerm),
		uuid:         uuid,
	}

	aa.init()

	return aa
}

// init sets default values for any unconfigured AchievementAssets properties.
func (aa *AchievementAssets) init() {
	if aa.maxFileSize == 0 {
		aa.maxFileSize = defaultMaxSize
	}

	if aa.imageQuality == 0 {
		aa.imageQuality = defaultQuality
	}

	if aa.dirPerm == 0 {
		aa.dirPerm = defaultDirPerm
	}

	if aa.filePerm == 0 {
		aa.filePerm = defaultFilePerm
	}
}

// getFileExt get file extension.
func (aa *AchievementAssets) getFileExt(filename string) string {
	return strings.TrimSuffix(filepath.Ext(filename), filepath.Base(filename))
}

// sanitizePath checks if the given path is safe to use.
// returns cleaned path or error if path is invalid.
func (aa *AchievementAssets) sanitizePath(path string) (string, error) {
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
func (aa *AchievementAssets) validateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > aa.maxFileSize {
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
func (aa *AchievementAssets) ensureUploadDirectory(uploadPath string) error {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrDirectoryMissing, uploadPath)
	}
	return nil
}

// readFileData reads file content with context support and size limitation.
func (aa *AchievementAssets) readFileData(ctx context.Context, file multipart.File) ([]byte, error) {
	done := make(chan struct{})
	var data []byte
	var readErr error

	go func() {
		data, readErr = io.ReadAll(io.LimitReader(file, aa.maxFileSize))
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
func (aa *AchievementAssets) convertToWebP(ctx context.Context, rawFile []byte) ([]byte, error) {
	done := make(chan struct{})
	var webp []byte
	var convertErr error

	go func() {
		options := bimg.Options{
			Quality: aa.imageQuality,
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
func (aa *AchievementAssets) UploadAndConvertToWebP(ctx context.Context, fileHeader *multipart.FileHeader) (achievementassets.UploadAndConvertToWebpResponse, error) {
	// check if context is already canceled.
	if err := ctx.Err(); err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, err
	}

	// measure upload duration.
	start := time.Now()
	defer func() {
		log.Printf("upload took %v", time.Since(start))
	}()

	// validate the uploaded file.
	if err := aa.validateFile(fileHeader); err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, err
	}

	// sanitize the target directory path.
	sanitizedDir, err := aa.sanitizePath(aa.dir)
	if err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("invalid directory: %w", err)
	}

	// open the uploaded file.
	file, err := fileHeader.Open()
	if err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	// prepare full upload path and verify directory exists.
	if err := aa.ensureUploadDirectory(sanitizedDir); err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, err
	}

	// read file data with context support.
	rawFile, err := aa.readFileData(ctx, file)
	if err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to read file: %w", err)
	}

	// convert image to WebP format.
	webp, err := aa.convertToWebP(ctx, rawFile)
	if err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to convert image: %w", err)
	}

	// generate unique filename and save the converted image.
	newName, err := aa.uuid.Generate()
	if err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to generate uuid: %w", err)
	}

	newFilePath := filepath.Join(sanitizedDir, newName+webpExt)
	if err := os.WriteFile(newFilePath, webp, aa.filePerm); err != nil {
		return achievementassets.UploadAndConvertToWebpResponse{}, fmt.Errorf("failed to save converted image: %w", err)
	}

	return achievementassets.UploadAndConvertToWebpResponse{
		Quality:        aa.imageQuality,
		NameFile:       newName + webpExt,
		ServerPathFile: filepath.Join(sanitizedDir, newName+webpExt),
		ClientPathFile: filepath.Join(aa.url, newName+webpExt),
		Extension:      webpExt,
		OldNameFile:    fileHeader.Filename,
		OldExtension:   aa.getFileExt(fileHeader.Filename),
	}, nil
}
