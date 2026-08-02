package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS short_long_mapping(
			short_url TEXT PRIMARY KEY,
			original_url TEXT
		);
	`)

	return nil

}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sql_db")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to sqlite3")
	return db, nil
}
