package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-jedi/lingvogramm_backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultQueryTimeout = 2000000000

// IPool defines the interface for the pgxpool.Pool.
//
//go:generate mockery --name=IPool --output=mocks --case=underscore
type IPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	Pool         IPool
	QueryTimeout int64

	host          string
	user          string
	password      string
	dbName        string
	port          int
	sslMode       string
	poolMaxConns  int
	migrationsDir string
}

func (p *Postgres) init() error {
	if p.QueryTimeout == 0 {
		p.QueryTimeout = defaultQueryTimeout
	}

	return nil
}

func New(ctx context.Context, cfg config.PostgresConfig) (*Postgres, error) {
	p := &Postgres{
		QueryTimeout:  cfg.QueryTimeout,
		host:          cfg.Host,
		user:          cfg.User,
		password:      cfg.Password,
		dbName:        cfg.DBName,
		port:          cfg.Port,
		sslMode:       cfg.SSLMode,
		poolMaxConns:  cfg.PoolMaxConns,
		migrationsDir: cfg.MigrationsDir,
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, p.generateDsn())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	p.Pool = pool

	if err := p.migrationsUp(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return p, nil
}

// generateDsn generate dsn string.
func (p *Postgres) generateDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s pool_max_conns=%d",
		p.host, p.user, p.password, p.dbName, p.port, p.sslMode, p.poolMaxConns,
	)
}

// migrationsUp up migrations for db.
func (p *Postgres) migrationsUp() error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.user, p.password, p.host, p.port, p.dbName, p.sslMode,
	)

	m, err := migrate.New(
		p.migrationsDir,
		dbURL,
	)
	if err != nil {
		return err
	}
	defer func(m *migrate.Migrate) {
		if err, _ := m.Close(); err != nil {
			log.Printf("error closes the source and the database: %v", err)
		}
	}(m)

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
