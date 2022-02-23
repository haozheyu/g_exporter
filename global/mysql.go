package global

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func InitConnect() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=10s", Option["mysql_user"],
		Option["mysql_password"], Option["mysql_host"], Option["mysql_port"], Option["mysql_database"]))
	fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=10s", Option["mysql_user"],
		Option["mysql_password"], Option["mysql_host"], Option["mysql_port"], Option["mysql_database"]))
	if err != nil {
		panic(fmt.Sprintln("Init mysql connect err,", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintln("Init mysql connect err,", err))
	}
	return db
}

func Connect(host, port, username, password, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=3s&readTimeout=5s", username, password, host, port, database))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}