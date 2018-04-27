package middleware

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Gorm 支持
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

//DefaultMySQL MySQL config
func DefaultMySQL() (*MySQL, error) {

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

	mysql := &MySQL{
		Gorm: db,
	}

	return mysql, nil
}

//DB *
type MySQL struct {
	Gorm  *gorm.DB
	Mutex sync.RWMutex
}

//MwMySQL MySQL middleware
func (m *MySQL) MwMySQL(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		c.Set("mysql", m.Gorm)
		return next(c)
	}
}
