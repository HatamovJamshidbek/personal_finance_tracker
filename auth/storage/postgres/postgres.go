package postgres

import (
	"auth_service/configs"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	_ "github.com/lib/pq"
)

type Store struct {
	db  *pgxpool.Pool
	cfg configs.Config
	log logger.ILogger
}

func (s Store) Users() storage.IUserStorage {
	return NewUserRepo(s.db, s.log)
}

func NewStore(ctx context.Context, cfg configs.Config, log logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("Failed to parse database config", logger.Error(err))
		return Store{}, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("Failed to create connection pool", logger.Error(err))
		return Store{}, err
	}

	//m, err := migrate.New("file://migrations/postgres", url)
	//if err != nil {
	//	log.Error("Failed to create migration instance", logger.Error(err))
	//	return Store{}, err
	//}
	//
	//if err = m.Up(); err != nil {
	//	if !strings.Contains(err.Error(), "no change") {
	//		log.Error("Migration failed", logger.Error(err))
	//		version, dirty, err := m.Version()
	//		if err != nil {
	//			log.Error("Failed to get migration version", logger.Error(err))
	//			return Store{}, err
	//		}
	//
	//		if dirty {
	//			version--
	//			if err = m.Force(int(version)); err != nil {
	//				log.Error("Failed to force migration version", logger.Error(err))
	//				return Store{}, err
	//			}
	//		}
	//		return Store{}, err
	//	}
	//}

	return Store{
		db:  pool,
		cfg: cfg,
		log: log,
	}, nil
}

func (s Store) Close() {
	s.db.Close()
}
