package data

import (
	"fmt"
	"log"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/dontpanicdao/caigo"
)

var (
	db         *sql.DB
	StarkCurve caigo.StarkCurve
)

func InitDB() {
	dbStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		Conf.JibeHost,
		5432,
		Conf.JibeUser,
		Conf.JibeDB,
		Conf.JibePassword)

	var err error
	db, err = sql.Open("postgres", dbStr)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
}

func CloseDB() {
	db.Close()
}

func InitStarkCuve() {
	var err error
	StarkCurve, err = caigo.SCWithConstants("./pedersen_params.json")
	if err != nil {
		log.Panic(err.Error())
	}
}
