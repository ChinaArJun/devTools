package global

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"oktools/src/conf"
)

var PSDB *sql.DB

func init() {
	dbConf := conf.Conf.Postgresdb
	var err error

	PSDB, err = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.DbName))
	if err != nil {
		log.Fatalln("Open SqlErr:", err)
	}
	log.Print("postgresdb connect success!")
}
