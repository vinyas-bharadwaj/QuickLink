package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS short_long_mapping (
			short_url VARCHAR(255) PRIMARY KEY,
			original_url TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func ConnectDB() (*sql.DB, error) {
	user := getEnv("MARIADB_USER", "appuser")
	password := getEnv("MARIADB_PASSWORD", "apppassword")
	host := getEnv("MARIADB_HOST", "localhost")
	port := getEnv("MARIADB_PORT", "3306")
	dbname := getEnv("MARIADB_DATABASE", "urlshortener")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to MariaDB: %w", err)
	}

	fmt.Println("Successfully connected to MariaDB")
	return db, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

