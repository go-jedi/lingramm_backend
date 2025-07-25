package achievementassets

import "os"

// constants for file server configuration.
const (
	defaultMaxSize  = 10 << 20          // default maximum file size (10MB)
	defaultQuality  = 30                // default image quality for conversion
	defaultDirPerm  = os.FileMode(0755) // default directory permissions
	defaultFilePerm = os.FileMode(0644) // default file permissions
	webpExt         = ".webp"           // extension for WebP images
)
