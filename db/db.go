package db

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server-center/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

func init() {
	dbConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: util.GormLog{ShowSql: config.Config.ShowSql, IgnoreErrs: []error{gorm.ErrRecordNotFound}},
	}

	var err error
	db, err = initMysql(dbConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(4)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(10)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func initMysql(dbConfig *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(config.Config.MysqlDsn), dbConfig)
}

func getDb(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func getTx(ctx context.Context) *gorm.DB {
	return getDb(ctx).Begin()
}
