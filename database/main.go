package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	return &Database{
		db: nil,
	}
}

func (d Database) Close() error {
	return d.db.Close()
}

func (d *Database) Initialize() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)

	db, err := sql.Open("sqlite3", exPath+"/scrapes.db")
	if err != nil {
		return err
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS scrapes (name TEXT PRIMARY KEY NOT NULL, result TEXT)")
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	d.db = db

	return nil
}

func (d Database) Update(name, result string) (bool, error) {
	statement, err := d.db.Prepare("SELECT result FROM scrapes WHERE name = ?")
	if err != nil {
		return false, err
	}

	row := statement.QueryRow(name)

	var storedValue string
	row.Scan(&storedValue)

	if storedValue == "" {
		statement, err = d.db.Prepare("INSERT INTO scrapes (name, result) VALUES (?, ?)")
		if err != nil {
			return false, err
		}
		_, err = statement.Exec(name, result)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if result != storedValue {
		statement, err = d.db.Prepare("UPDATE scrapes SET result = ? WHERE name = ?")
		if err != nil {
			return false, err
		}
		_, err = statement.Exec(result, name)
		if err != nil {
			return false, err
		}

		return true, nil
	}
	return false, nil
}
