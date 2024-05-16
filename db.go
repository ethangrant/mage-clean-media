package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func DbConnect(user string, password string, host string, dbName string) (db *sql.DB, err error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		return nil, errors.New("There was a problem connecting to the database: " + err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, errors.New("Could not ping database: " + err.Error())
	}

	return db, nil
}