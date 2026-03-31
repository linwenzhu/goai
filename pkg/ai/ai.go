package ai

import (
	"github.com/sashabaranov/go-openai"
	"goAi/config"
)

var Client *openai.Client

func Init(cfg config.AIConfig) {
	clientCfg := openai.DefaultConfig(cfg.APIKey)
	clientCfg.BaseURL = cfg.BaseURL
	Client = openai.NewClientWithConfig(clientCfg)
}
