package database

import (
	"database/sql"
	"fmt"
	"traffic-guarder/internal/infrastructure/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func New(cfg config.DBConfig) *bun.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		panic(err)
	}

	AutoMigration(db)

	return db
}
