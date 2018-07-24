// mysql-helpers
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type OneRow struct {
	id     int
	a_id   int
	b_id   int
	action string
	status string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "microservice"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Println(err)
	}
	return db
}

func showAllRows() {
	db := dbConn()

	query, err := db.Query("SELECT * FROM microservice.actions")
	if err != nil {
		log.Println(err)
	}
	defer query.Close()

	rows := make([]*OneRow, 0)
	for query.Next() {
		row := new(OneRow)
		err := query.Scan(&row.id, &row.a_id, &row.b_id, &row.action, &row.status)
		if err != nil {
			log.Println(err)
		}
		rows = append(rows, row)
	}
	if err = query.Err(); err != nil {
		log.Println(err)
	}

	for _, row := range rows {
		fmt.Printf("%d, %d, %d, %s, %s\n", row.id, row.a_id, row.b_id, row.action, row.status)
	}
	defer db.Close()
}

func findAId(newAId int, act string) (flag bool) {
	db := dbConn()

	query, err := db.Query("SELECT * FROM microservice.actions")
	if err != nil {
		log.Println(err)
	}
	defer query.Close()

	rows := make([]*OneRow, 0)
	for query.Next() {
		row := new(OneRow)
		err := query.Scan(&row.id, &row.a_id, &row.b_id, &row.action, &row.status)
		if err != nil {
			log.Println(err)
		}
		rows = append(rows, row)
	}
	if err = query.Err(); err != nil {
		log.Println(err)
	}

	for _, row := range rows {
		if row.a_id == newAId && row.action == act {
			return true
		}
	}
	defer db.Close()
	return false
}

func updateBId(newBId int, act string, errorB bool) (flag bool) {
	db := dbConn()

	query, err := db.Query("SELECT * FROM microservice.actions")
	if err != nil {
		log.Println(err)
	}
	defer query.Close()

	rows := make([]*OneRow, 0)
	for query.Next() {
		row := new(OneRow)
		err := query.Scan(&row.id, &row.a_id, &row.b_id, &row.action, &row.status)
		if err != nil {
			log.Println(err)
		}
		rows = append(rows, row)
	}
	if err = query.Err(); err != nil {
		log.Println(err)
	}

	for _, row := range rows {
		if row.b_id == 0 && row.action == act {
			updDb, err := db.Prepare("UPDATE microservice.actions SET b_operation_id=?, status=? WHERE b_operation_id=? and action=? and status=?")
			if err != nil {
				log.Println(err)
			}
			if errorB == false {
				updDb.Exec(newBId, "success", 0, act, "in_progress")
			} else {
				updDb.Exec(0, "error", 0, act, "in_progress")
			}

			defer updDb.Close()
			return true
		}
	}
	defer db.Close()
	return false
}

func insertNewAction(newId int, act string) {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO microservice.actions(a_operation_id, b_operation_id, action, status) VALUES(?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	insForm.Exec(newId, 0, act, "in_progress")
	log.Println("INSERT: a_operation_id: " + strconv.Itoa(newId) + " | action: " + act)
	defer db.Close()
}
