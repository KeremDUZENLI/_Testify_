package database

import (
	"testify/common/env"
	"testify/common/helper"

	"database/sql"
	"fmt"
	"time"
)

var postgreDB *sql.DB

func DbConnect(args ...string) *sql.DB {
	var databaseName string
	if args == nil {
		databaseName = env.DbName
	} else {
		databaseName = args[0]
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.DbHost, env.DbPort, env.DbUser, env.DbPass, databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	helper.ErrorLog(err)
	helper.ErrorLog(db.Ping())

	postgreDB = db
	return postgreDB
}

func DbCreateTable() {
	stmt, err := postgreDB.Prepare(
		`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp TIMESTAMPTZ PRIMARY KEY,
		price DECIMAL NOT NULL
	)`)
	helper.ErrorLog(err)

	_, err = stmt.Exec()
	helper.ErrorLog(err)
}

func DbSeedTable() {
	var rowCount int
	postgreDB.QueryRow("SELECT COUNT(*) FROM stockprices").Scan(&rowCount)

	if rowCount != 5 {
		for i := 1; i <= 5; i++ {
			_, err := postgreDB.Exec("INSERT INTO stockprices (timestamp, price) VALUES ($1,$2)",
				time.Now().Add(time.Duration(-i)*time.Minute), float64((6-i)*5))

			helper.ErrorLog(err)
		}
	}
}

func DbCreateExtra() {
	_, err := postgreDB.Exec(`CREATE DATABASE postgres_test`)
	helper.ErrorPrint(err)
}
