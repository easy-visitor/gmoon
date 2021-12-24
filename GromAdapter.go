package gmoon

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type GromAdapter struct {
	*gorm.DB
}

func dns() string {
	return config.Mysql.Username + ":" + config.Mysql.Password + "@tcp(" + config.Mysql.Path + ":" + config.Mysql.Port + ")/" + config.Mysql.Dbname + "?" + config.Mysql.Config
}

func NewGromAdapter() *GromAdapter {
	mysqlConfig := mysql.Config{
		DSN:                       dns(), // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), dbConfig())
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	return &GromAdapter{
		DB: db,
	}
}

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

func dbConfig() *gorm.Config {
	gormConfig := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	switch config.Mysql.LogMode {
	case "silent", "Silent":
		gormConfig.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		gormConfig.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		gormConfig.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		gormConfig.Logger = _default.LogMode(logger.Info)
	default:
		gormConfig.Logger = _default.LogMode(logger.Info)
	}
	return gormConfig
}
