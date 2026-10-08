package database

import (
	"database/sql"
	"embed"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed migrations.sql
var migrationsFS embed.FS

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite", "file:SOKOLOV_BOTS.db?cache=shared")
	if err != nil {
		return err
	}

	migrations, err := migrationsFS.ReadFile("migrations.sql")
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(migrations))
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
