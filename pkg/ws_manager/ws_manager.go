package wsmanager

import notificationhub "github.com/go-jedi/lingramm_backend/pkg/ws_manager/notification"

type WSManager struct {
	NotificationHUB *notificationhub.Hub
}

func New() *WSManager {
	return &WSManager{
		NotificationHUB: notificationhub.New(),
	}
}
