package db

import (
	"fmt"
	"scaffold/pkg/config"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(cfg *config.DBConfig) error {
	var dialector gorm.Dialector

	// 根据驱动类型选择数据库驱动
	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host,
			cfg.Username,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
		)
		dialector = postgres.Open(dsn)

	default:
		return errors.New(fmt.Sprintf("不支持的数据库驱动: %s", cfg.Driver))
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg.LogLevel)),
	}

	// 建立连接
	var err error
	db, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return errors.Wrap(err, "数据库连接失败")
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "获取底层sqlDB失败")
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenCon)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleCon)
	sqlDB.SetConnMaxLifetime(time.Hour)

	zap.L().Info("数据库连接成功",
		zap.String("driver", cfg.Driver),
		zap.String("host", cfg.Host),
		zap.String("database", cfg.DBName))

	return nil
}

func Close() {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		zap.L().Error("db 关闭连接失败", zap.Error(err))
	}
}

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

func DB() *gorm.DB {
	return db
}

// func Migrate() error {
// 	if err := db.AutoMigrate("xxx"); err != nil {
// 		zap.L().Error("数据库迁移失败", zap.Error(err))
// 		return err
// 	}
// 	zap.L().Info("数据库迁移成功")
// 	return nil
// }
