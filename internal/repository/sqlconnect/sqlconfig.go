package sqlconnect

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	db, err := sql.Open("sqlite3", "./sql_db")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
