package db

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Store *gorm.DB

func init() {
	conf := mysqlDriver.NewConfig()
	conf.Addr = "127.0.0.1:3306"
	conf.User = "ocean"
	conf.Passwd = "oceanho"
	conf.DBName = "djdb"
	Store, err := gorm.Open(mysql.Open(conf.FormatDSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("initial db. %v", err))
	}
	db, err := Store.DB()
	if err != nil {
		panic(fmt.Errorf("get sql.db. %v", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping db. %v", err))
	}
}
