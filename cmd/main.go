package main

import (
	"fmt"
	"github.com/spf13/viper"
	"goAi/config"
	"goAi/internal/router"
	"goAi/pkg/ai"
	"goAi/pkg/database"
	"goAi/pkg/jwt"
	"goAi/pkg/logger"
)

var Cfg config.Config

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic("解析配置文件失败: " + err.Error())
	}
}
func main() {
	// 1. 加载配置
	initConfig()

	// 2. 初始化日志
	logger.Init(Cfg.Server.Mode)
	defer logger.Log.Sync()

	// 3. 初始化 MySQL
	database.InitMySQL(Cfg.Database)

	// 4. 初始化 Redis
	database.InitRedis(Cfg.Redis)

	// 5. 初始化 JWT
	jwt.Init(Cfg.JWT)

	ai.Init(Cfg.AI)
	logger.Info("AI 初始化成功")
	// 6. 启动服务
	r := router.NewRouter()

	addr := fmt.Sprintf(":%d", Cfg.Server.Port)
	logger.Info("服务启动，监听地址: " + addr)
	r.Run(addr)
}
