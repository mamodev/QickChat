package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init(path string) error {

	var err error
	db, err = sql.Open("sqlite3", path)
	
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(10)

	return err
}

func Close() {
	db.Close()
}

func Conn(ctx context.Context) (*sql.Conn, error) {
	conn, err := db.Conn(ctx)
	return conn, err
}
