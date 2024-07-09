package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Enotisi/go_final_project/internal/config"
	_ "modernc.org/sqlite"
)

const (
	baseName = "scheduler.db"
)

func MustBeInitDB() *sql.DB {

	install := checkDataBaseFile()

	db, err := sql.Open("sqlite", baseName)

	if err != nil {
		panic(err.Error())
	}

	if !install {
		createTableAndIndex(db)
	}

	return db
}

func checkDataBaseFile() bool {

	path := config.Conf.BasePath
	if path == "" {
		appPath, err := os.Getwd()

		if err != nil {
			panic(err.Error())
		}
		path = appPath
	}

	basePath := filepath.Join(path, baseName)

	_, err := os.Stat(basePath)

	return err == nil
}

func createTableAndIndex(db *sql.DB) {

	createTableSql := `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL CHECK(length(date) == 8),
	title TEXT NOT NULL,
	comment TEXT,
	repeat TEXT CHECK(length(repeat) <= 128)
	);`
	_, err := db.Exec(createTableSql)

	if err != nil {
		panic(err.Error())
	}

	createIndexSql := `CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);`

	_, err = db.Exec(createIndexSql)

	if err != nil {
		panic(err.Error())
	}

}
