package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	mysql "github.com/go-sql-driver/mysql"
)

func Init() *sql.DB {
	mysqlHost, provided := os.LookupEnv("MYSQL_HOST")
	if !provided {
		mysqlHost = "localhost"
	}
	mysqlUser, provided := os.LookupEnv("MYSQL_USER")
	if !provided {
		log.Fatalf("Environment variable %s is not set", "MYSQL_USER")
	}
	mysqlPass, provided := os.LookupEnv("MYSQL_PASS")
	if !provided {
		log.Fatalf("Environment variable %s is not set", "MYSQL_PASS")
	}
	config := mysql.Config{
		User:                 mysqlUser,
		Passwd:               mysqlPass,
		Net:                  "tcp",
		Addr:                 mysqlHost,
		AllowNativePasswords: true,
		DBName:               "placesdb",
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	checkErr(err)

	err = db.Ping()
	checkErr(err)
	fmt.Printf("Connection successfully")

	return db
}

func checkErr(err error) {
	if err != nil {
		fmt.Print(err.Error())
	}
}
