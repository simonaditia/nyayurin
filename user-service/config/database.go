package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/simonaditia/nyayurin/user-service/database/seeds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.Dbname)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] Failed to connect to database" + cfg.Psql.Host)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] Failed to get database instance")
		return nil, err
	}

	seeds.SeedRole(db)
	seeds.SeedAdmin(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
