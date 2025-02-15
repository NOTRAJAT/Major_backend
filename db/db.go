package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMYSQLStorage(cfg mysql.Config) (*sql.DB,error){
	
	db,err:= sql.Open("mysql",cfg.FormatDSN())
	if err!=nil{
		log.Fatal(err)
	}
	return db,nil 

}

func InitStorage(db *sql.DB){
	if err:= db.Ping(); err!=nil{
		 log.Fatal(err)
	}
	log.Println("Db Connected...")
}