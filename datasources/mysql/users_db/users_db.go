package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const (
	my_users_db_username = "my_users_db_username"
	my_users_db_password = "my_users_db_password"
	my_users_db_host     = "my_users_db_host"
	my_users_db_scheme   = "my_users_db_scheme"
)

var (
	Client *sql.DB

	username = os.Getenv(my_users_db_username)
	password = os.Getenv(my_users_db_password)
	host     = os.Getenv(my_users_db_host)
	scheme   = os.Getenv(my_users_db_scheme)
)

func init() {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, scheme,
	)
	var err error
	Client, err = sql.Open("mysql", datasourceName)
	if err != nil {
		panic(err)
	}
	if err := Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfully configured")
}
