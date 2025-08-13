package subscription

import (
	createsubscription "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription/create_subscription"
	createsubscriptionhistory "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription/create_subscription_history"
	existsbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription/exists_by_telegram_id"
	getbytelegramid "github.com/go-jedi/lingramm_backend/internal/repository/v1/subscription/get_by_telegram_id"
	"github.com/go-jedi/lingramm_backend/pkg/logger"
)

type Repository struct {
	CreateSubscription        createsubscription.ICreateSubscription
	CreateSubscriptionHistory createsubscriptionhistory.ICreateSubscriptionHistory
	ExistsByTelegramID        existsbytelegramid.IExistsByTelegramID
	GetByTelegramID           getbytelegramid.IGetByTelegramID
}

func New(
	queryTimeout int64,
	logger logger.ILogger,
) *Repository {
	return &Repository{
		CreateSubscription:        createsubscription.New(queryTimeout, logger),
		CreateSubscriptionHistory: createsubscriptionhistory.New(queryTimeout, logger),
		ExistsByTelegramID:        existsbytelegramid.New(queryTimeout, logger),
		GetByTelegramID:           getbytelegramid.New(queryTimeout, logger),
	}
}
