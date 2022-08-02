package starter

import (
	"database/sql"
	"github.com/XM-GO/PandaKit/config"
	"github.com/XM-GO/PandaKit/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"

	gormlog "gorm.io/gorm/logger"

	_ "github.com/lib/pq"
)

var Db *gorm.DB

func GormInit(ty string) *gorm.DB {
	switch ty {
	case "mysql":
		Db = GormMysql()
	case "postgresql":
		Db = GormPostgresql()
	}
	return Db
}
func GormMysql() *gorm.DB {
	m := config.Conf.Mysql
	if m == nil || m.Dbname == "" {
		logger.Log.Panic("未找到数据库配置信息")
		return nil
	}
	logger.Log.Infof("连接mysql [%s]", m.Host)
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	ormConfig := &gorm.Config{Logger: gormlog.Default.LogMode(gormlog.Silent)}
	db, err := gorm.Open(mysql.New(mysqlConfig), ormConfig)
	if err != nil {
		logger.Log.Panicf("连接mysql失败! [%s]", err.Error())
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}

func GormPostgresql() *gorm.DB {
	m := config.Conf.Postgresql
	if m == nil || m.Dbname == "" {
		logger.Log.Panic("未找到数据库配置信息")
		return nil
	}
	logger.Log.Infof("连接postgres [%s]", m.PgDsn())
	db, err := sql.Open("postgres", m.PgDsn())
	if err != nil {
		logger.Log.Panicf("连接postgresql失败! [%s]", err.Error())
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		logger.Log.Panicf("连接postgresql失败! [%s]", err.Error())
	}
	sqlDB, err := gormDb.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDb
}
