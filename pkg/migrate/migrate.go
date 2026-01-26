package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../../db/migrations")

	// Open a separate connection for migrations using pgx stdlib
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Panicf("[MIGRATION] failed to connect: %v", err)
	}
	defer db.Close()

	// Set PostgreSQL dialect for Goose
	if err := goose.SetDialect("postgres"); err != nil {
		log.Panicf("[MIGRATION] set dialect error: %v", err)
	}

	fmt.Println("[MIGRATION] running...")

	if err := goose.Up(db, rootPath); err != nil {
		log.Panicf("[MIGRATION] up error: %v", err)
	}

	fmt.Println("[MIGRATION] done...")
}
