package model

import (
	"database/sql"
	"fmt"

	"github.com/Bit0r/online-store/conf"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var config struct {
		User      string
		Password  string
		Host      string
		Port      uint16
		Name      string
		ParseTime bool
	}

	conf.Unmarshal("database", "online_store_database", &config)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=%v",
		config.User, config.Password,
		config.Host, config.Port,
		config.Name, config.ParseTime)
	db, _ = sql.Open("mysql", dsn)
}
