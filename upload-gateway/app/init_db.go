package app

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var (
	db *gorm.DB
)

func GetDb() *gorm.DB {
	if db == nil {
		panic("database connection is not initialized")
	}

	return db
}

func InitDb() {
	var (
		connection *gorm.DB
		sqlDb      *sql.DB
		err        error
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=UTC&timeout=1m&readTimeout=1m",
		"homestead",
		"secret",
		"mysql",
		3306,
		"homestead",
	)

	if connection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: false,
			TablePrefix:   "",
		},
	}); err != nil {
		panic(err)
	}

	if sqlDb, err = connection.DB(); err != nil {
		panic(err)
	}

	sqlDb.SetConnMaxLifetime(time.Minute * 1)
	sqlDb.SetConnMaxIdleTime(time.Minute * 1)
	sqlDb.SetMaxIdleConns(50)
	sqlDb.SetMaxOpenConns(100)

	db = connection.WithContext(log.Logger.WithContext(context.Background()))
}
