package fileserver

import (
	"github.com/go-jedi/lingramm_backend/config"
	achievementassets "github.com/go-jedi/lingramm_backend/pkg/file_server/achievement_assets"
	awardassets "github.com/go-jedi/lingramm_backend/pkg/file_server/award_assets"
	clientassets "github.com/go-jedi/lingramm_backend/pkg/file_server/client_assets"
	"github.com/go-jedi/lingramm_backend/pkg/uuid"
)

type FileServer struct {
	AchievementAssets achievementassets.IAchievementAssets
	AwardAssets       awardassets.IAwardAssets
	ClientAssets      clientassets.IClientAssets
}

func New(cfg config.FileServerConfig, uuid *uuid.UUID) *FileServer {
	return &FileServer{
		AchievementAssets: achievementassets.New(cfg, uuid),
		AwardAssets:       awardassets.New(cfg, uuid),
		ClientAssets:      clientassets.New(cfg, uuid),
	}
}
