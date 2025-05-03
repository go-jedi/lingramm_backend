package signin

import (
	"context"
	"log"

	"github.com/go-jedi/lingvogramm_backend/internal/domain/auth"
	"github.com/go-jedi/lingvogramm_backend/internal/domain/user"
	userrepository "github.com/go-jedi/lingvogramm_backend/internal/repository/user"
	"github.com/go-jedi/lingvogramm_backend/pkg/logger"
	"github.com/go-jedi/lingvogramm_backend/pkg/postgres"
	"github.com/go-jedi/lingvogramm_backend/pkg/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockery --name=ISignIn --output=mocks --case=underscore
type ISignIn interface {
	Execute(ctx context.Context, dto auth.SignInDTO) (user.User, error)
}

type SignIn struct {
	userRepository *userrepository.Repository
	logger         logger.ILogger
	postgres       *postgres.Postgres
	uuid           uuid.IUUID
}

func New(
	userRepository *userrepository.Repository,
	logger logger.ILogger,
	postgres *postgres.Postgres,
	uuid uuid.IUUID,
) *SignIn {
	return &SignIn{
		userRepository: userRepository,
		logger:         logger,
		postgres:       postgres,
		uuid:           uuid,
	}
}

func (si *SignIn) Execute(ctx context.Context, dto auth.SignInDTO) (user.User, error) {
	si.logger.Debug("[sign in user] execute service")

	var (
		err error
		u   user.User
	)

	tx, err := si.postgres.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return user.User{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("failed rollback the transaction: %v", rbErr)
			}
		}
	}()

	ie, err := si.userRepository.Exists.Execute(ctx, tx, dto.TelegramID, dto.Username)
	if err != nil {
		return user.User{}, err
	}

	if ie {
		u, err = si.findOrReturnExisting(ctx, tx, dto.TelegramID)
	} else {
		u, err = si.createUser(ctx, tx, dto)
	}
	if err != nil {
		return user.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// createUser create new user.
func (si *SignIn) createUser(ctx context.Context, tx pgx.Tx, dto auth.SignInDTO) (user.User, error) {
	newUUID, err := si.uuid.Generate()
	if err != nil {
		return user.User{}, err
	}

	createDTO := user.CreateDTO{
		UUID:       newUUID,
		TelegramID: dto.TelegramID,
		Username:   dto.Username,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
	}

	return si.userRepository.Create.Execute(ctx, tx, createDTO)
}

// findOrReturnExisting find user if user existing.
func (si *SignIn) findOrReturnExisting(ctx context.Context, tx pgx.Tx, telegramID string) (user.User, error) {
	return si.userRepository.GetByTelegramID.Execute(ctx, tx, telegramID)
}
