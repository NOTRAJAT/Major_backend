package main

import (
	"database/sql"
	"fmt"
	"log"
	"myAttendance/cmd/api"
	"myAttendance/config"
	"myAttendance/db"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db,err:= db.NewMYSQLStorage(mysql.Config{
		User: config.ENV.DbUser,
		Passwd: config.ENV.DbPassword,
		Addr: config.ENV.DbAddress,
		DBName: config.ENV.DbName,
		Net: "tcp",	
		AllowNativePasswords: true,
		ParseTime: true,
	})

	if err!=nil{
		 log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf("%s:%s",config.ENV.PublicHost,config.ENV.Port),db) // localhost:3000
	if err:=server.Run();err!=nil{
		log.Fatal(err)
	}

}


func initStorage(db *sql.DB){
	if err:= db.Ping(); err!=nil{
		 log.Fatal(err)
	}
	log.Println("Db Connected...")
}
