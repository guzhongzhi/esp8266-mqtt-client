package tv

import (
	"database/sql"
	"fmt"
)

var dbFileName = ""

var connection *sql.DB

func SetDbFileName(fileName string) error {
	dbFileName = fileName
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	connection = db
	return nil
}

func CreateTables() error {
	sql := `
	CREATE TABLE IF NOT EXISTS "ir_received" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"mac" VARCHAR(64) NOT NULL default "",
		"data" TEXT NOT NULL default "",
		"device" VARCHAR(64) NOT NULL default "",
		"name" VARCHAR(64) NOT NULL default "",
		"created_at" TIMESTAMP default (datetime('now', 'localtime')),
		"updated_at" TIMESTAMP default (datetime('now', 'localtime')),
	PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS "user" (
	"id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"username" VARCHAR(64) NOT NULL default "",
	"mac" VARCHAR(64) NOT NULL default "",
	"ip"  VARCHAR(64) NOT NULL default "",
	"wifi"  VARCHAR(64) NOT NULL default "",
	"relay"  VARCHAR(64) NOT NULL default "", 
	"created_at" TIMESTAMP default (datetime('now', 'localtime')),
	"updated_at" TIMESTAMP default (datetime('now', 'localtime')),
	PRIMARY KEY (id)
);
`
	rs,err := connection.Exec(sql)
	if err != nil {
		return err
	}
	fmt.Println("rs",rs)
	return nil
}
