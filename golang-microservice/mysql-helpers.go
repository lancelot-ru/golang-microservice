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
	checkError(err, "opening DB")
	return db
}

func showAllRows() {
	db := dbConn()

	query, err := db.Query("SELECT * FROM microservice.actions")
	checkError(err, "selecting from DB")
	defer query.Close()

	rows := make([]*OneRow, 0)
	for query.Next() {
		row := new(OneRow)
		err := query.Scan(&row.id, &row.a_id, &row.b_id, &row.action, &row.status)
		checkError(err, "scaning DB")
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

func findAId(newAId int, act string) (flag int) {
	db := dbConn()

	var a_operation_id int
	var action string

	sqlString := "SELECT a_operation_id, action FROM microservice.actions WHERE a_operation_id=" + strconv.Itoa(newAId) + " AND action='" + act + "';"

	row := db.QueryRow(sqlString)
	switch err := row.Scan(&a_operation_id, &action); err {
	case sql.ErrNoRows:
		return 0
	case nil:
		return 1
	default:
		checkError(err, "selecting from DB")
	}

	defer db.Close()
	return -1
}

func insertNewAction(newId int, act string) {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO microservice.actions(a_operation_id, b_operation_id, action, status) VALUES(?,?,?,?)")
	checkError(err, "preparing DB")
	insForm.Exec(newId, 0, act, "in_progress")
	log.Println("INSERT: a_operation_id: " + strconv.Itoa(newId) + " | action: " + act)
	defer insForm.Close()
	defer db.Close()
}

func updateBId(newBId int, act string, errorB bool) (flag bool) {
	db := dbConn()
	defer db.Close()

	updDb, err := db.Prepare("UPDATE microservice.actions SET b_operation_id=?, status=? WHERE b_operation_id=? and action=? and status=?")
	defer updDb.Close()

	checkError(err, "preparing DB")
	if errorB == false {
		updDb.Exec(newBId, "success", 0, act, "in_progress")
		return true
	} else {
		updDb.Exec(0, "error", 0, act, "in_progress")
		return false
	}

}
