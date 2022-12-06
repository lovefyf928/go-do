package mysql

import (
	"go-do/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var MysqlDb *gorm.DB

func init() {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.ConfigInfo.DataSource.Mysql.Uri,                       // DSN data source name
		DefaultStringSize:         conf.ConfigInfo.DataSource.Mysql.DefaultStringSize,         // string 类型字段的默认长度
		DisableDatetimePrecision:  conf.ConfigInfo.DataSource.Mysql.DisableDatetimePrecision,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    conf.ConfigInfo.DataSource.Mysql.DontSupportRenameIndex,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   conf.ConfigInfo.DataSource.Mysql.DontSupportRenameColumn,   // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: conf.ConfigInfo.DataSource.Mysql.SkipInitializeWithVersion, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(conf.ConfigInfo.DataSource.Mysql.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(conf.ConfigInfo.DataSource.Mysql.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(conf.ConfigInfo.DataSource.Mysql.ConnMaxLifetime))

	MysqlDb = db

}
