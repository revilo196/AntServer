package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

const DbFILENAME = "./data.db"

const creatStmt = `create table data (timestp integer not null primary key, temperature real);`
const insertStmt = `insert into data(timestp, temperature) values(?, ?)`
const checkStmt = "select name from sqlite_master WHERE type='table'"
const selectStmt = "select * from data where ? < timestp"

var db *sql.DB = nil

func ResetDB() {
	CloseDB()
	err := os.Remove(DbFILENAME)
	if err != nil {
		log.Fatal(err)
	}
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", DbFILENAME)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckBaseDB() bool {
	if db == nil {
		log.Println("DB not Init")
		return false
	}
	row, err := db.Query(checkStmt)
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		var name string
		err = row.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		if name == "data" {
			return true
		}
	}

	return false
}

func BuildBaseDB() {
	if db == nil {
		log.Println("DB not Init")
		return
	}
	_, err := db.Exec(creatStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt)
		return
	}
}

func CloseDB() {
	if db == nil {
		log.Println("DB not Init")
		return
	}
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func AddValuesToDB(temperature []float32, times []time.Time) {
	if db == nil {
		log.Println("DB not Init")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(insertStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < len(temperature) || i < len(times); i++ {
		_, err = stmt.Exec(times[i].Unix(), temperature[i])
		if err != nil {
			log.Println(tx.Rollback())
			log.Fatal(err)
		}
		fmt.Println("Recived", times[i], temperature[i])
	}

	err = tx.Commit()
	if err != nil {
		log.Println(tx.Rollback())
		log.Fatal(err)
	}
}

func GetValuesFromDB(since time.Time) ([]float32, []time.Time) {
	temperature := make([]float32, 0)
	timestamps := make([]time.Time, 0)

	if db == nil {
		log.Println("DB not Init")
		return temperature, timestamps
	}

	fmt.Println(since.Unix())

	row, err := db.Query(selectStmt, since.Unix())
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		var timestp int64
		var temp float32
		err = row.Scan(&timestp, &temp)
		if err != nil {
			log.Fatal(err)
		}
		temperature = append(temperature, temp)
		timestamps = append(timestamps, time.Unix(timestp, 0))
	}
	return temperature, timestamps
}
