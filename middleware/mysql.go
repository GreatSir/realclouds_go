package middleware

import (
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Gorm 支持
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

//MySQLConf MySQL config
func MySQLConf() (*gorm.DB, error) {

	devMode := utils.GetENVToBool("DEV_MODE")

	dbHost := utils.GetENV("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "127.0.0.1:3306"
	}

	dbUserName := utils.GetENV("DB_USERNAME")
	if len(dbUserName) == 0 {
		dbUserName = "ayerdudu"
	}

	dbPassword := utils.GetENV("DB_PASSWORD")
	if len(dbPassword) == 0 {
		dbPassword = "AyerDudu888"
	}

	dbDataBase := utils.GetENV("DB_DATABASE")
	if len(dbDataBase) == 0 {
		dbDataBase = "ayerdudu"
	}

	dbMaxIdleConns, err := utils.GetENVToInt("DB_MAXIDLECONNS")
	if nil != err {
		dbMaxIdleConns = 10
	}

	dbMaxOpenConns, err := utils.GetENVToInt("DB_MAXOPENCONNS")
	if nil != err {
		dbMaxOpenConns = 100
	}

	dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Asia%%2FShanghai&timeout=30s", dbUserName, dbPassword, dbHost, dbDataBase)

	db, err := gorm.Open("mysql", dbURL)
	if nil != err {
		return nil, err
	}

	db.DB().SetMaxIdleConns(dbMaxIdleConns)
	db.DB().SetMaxOpenConns(dbMaxOpenConns)

	db.LogMode(devMode)

	if err = db.DB().Ping(); nil != err {
		return nil, err
	}

	return db, nil
}

//DB *
type DB struct {
	Gorm  *gorm.DB
	Mutex sync.RWMutex
}

//MySQL New MySQL dirver
func MySQL() *DB {
	db, err := MySQLConf()
	if nil != err {
		log.Fatalf("New mysql driver error: %v", err.Error())
		return nil
	}
	return &DB{Gorm: db}
}

//MwMySQL MySQL middleware
func (d *DB) MwMySQL(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		d.Mutex.Lock()
		defer d.Mutex.Unlock()
		c.Set("mysql", d.Gorm)
		return next(c)
	}
}
