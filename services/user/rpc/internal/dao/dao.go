package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
	"wonderful-hand-user/rpc/internal/config"
)

var DB *gorm.DB

func init() {
	c, _ := config.ReadConfig()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Mysql.DataSourceName,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalln(err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(c.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(c.Mysql.ConnMaxLifetime))
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
}
