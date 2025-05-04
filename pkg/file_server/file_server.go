package fileserver

import (
	"github.com/go-jedi/lingvogramm_backend/config"
	clientassets "github.com/go-jedi/lingvogramm_backend/pkg/file_server/client_assets"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
)

type FileServer struct {
	ClientAssets clientassets.IClientAssets
}

func New(cfg config.FileServerConfig, uuid *uuid.UUID) *FileServer {
	return &FileServer{
		ClientAssets: clientassets.New(cfg, uuid),
	}
}
