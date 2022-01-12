package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dontpanicdao/caigo"
	_ "github.com/lib/pq"
)

var (
	db         *sql.DB
	StarkCurve caigo.StarkCurve
	TypedCert  caigo.TypedData
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

func InitTypedCert(chainId int) (err error) {
	certTypes := make(map[string]caigo.TypeDef)
	snDefs := []caigo.Definition{
		caigo.Definition{"name", "felt"},
		caigo.Definition{"version", "felt"},
		caigo.Definition{"chainId", "felt"},
	}
	certTypes["StarkNetDomain"] = caigo.TypeDef{Definitions: snDefs}

	certDefs := []caigo.Definition{
		caigo.Definition{"certUri", "felt"},
		caigo.Definition{"certKey", "felt"},
	}

	certTypes["Cert"] = caigo.TypeDef{Definitions: certDefs}
	dm := caigo.Domain{
		Name:    "StarkNet Cert",
		Version: "1",
		ChainId: chainId,
	}

	TypedCert, err = caigo.NewTypedData(certTypes, "Cert", dm)
	return err
}
