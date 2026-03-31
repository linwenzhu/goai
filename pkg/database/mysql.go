package database

import (
	"fmt"
	"goAi/config"
	"goAi/internal/model"
	"goAi/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL(cfg config.DatabaseConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("MySQL 连接失败: " + err.Error())
	}

	// 自动迁移，有表不动，没表自动创建
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		logger.Fatal("数据库迁移失败: " + err.Error())
	}

	logger.Info("MySQL 连接成功")
}
