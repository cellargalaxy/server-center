package db

import (
	"context"
	"github.com/cellargalaxy/go_common/util"
	"github.com/cellargalaxy/server_center/config"
	"github.com/cellargalaxy/server_center/model"
	"github.com/glebarez/sqlite"
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
	db, err = initDb(dbConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(1)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(1)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&model.ServerConfModel{})
	if err != nil {
		panic(err)
	}
}

func initDb(dbConfig *gorm.Config) (*gorm.DB, error) {
	if config.Config.MysqlDsn != "" {
		return gorm.Open(mysql.Open(config.Config.MysqlDsn), dbConfig)
	}
	err := util.CreateFolderPath(util.CreateLogCtx(), "resource")
	if err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open("resource/sqlite.db"), dbConfig)
}

func getDb(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func getTx(ctx context.Context) *gorm.DB {
	return getDb(ctx).Begin()
}
