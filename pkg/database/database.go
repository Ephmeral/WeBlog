package database

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB对象
var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector, _logger gormlogger.Interface) {
	// 使用gorm.Open连接数据库
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{Logger: _logger})
	// 处理错误
	if err != nil {
		fmt.Println(err.Error())
	}

	// 获取底层的sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}
