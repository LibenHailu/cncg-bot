package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Telegram struct {
		BotToken  string
		ChannelID string
		ParseMode string
	}
	Scheduler struct {
		CronSpec  string
		BatchSize int
	}
	Filters struct {
		MaxAgeDays int
		MinScore   float64
	}
	Keywords struct {
		Positive []string
		Negative []string
	}
	Sources []struct {
		Name string
		Type string
		URL  string
		Weight float64
		Tags []string
	}
	DBPath string
}

func Load(path string) (Config, error) {
	var cfg Config
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil { return cfg, err }
	if err := v.Unmarshal(&cfg); err != nil { return cfg, err }
	if cfg.DBPath == "" { cfg.DBPath = "data.db" }
	// TODO: This should be loaded from the config.yaml
	if cfg.Telegram.ParseMode == "" { cfg.Telegram.ParseMode = "MarkdownV2" }
	if cfg.Telegram.ChannelID == "" { cfg.Telegram.ChannelID = "@CloudNativeAddisAbaba" }
	if cfg.Telegram.BotToken == "" { cfg.Telegram.BotToken = "" }
	if cfg.Scheduler.CronSpec == "" { cfg.Scheduler.CronSpec = "* * * * *" }
	if cfg.Scheduler.BatchSize == 0 { cfg.Scheduler.BatchSize = 4 }
	if cfg.Filters.MaxAgeDays == 0 { cfg.Filters.MaxAgeDays = 1000 }
	if cfg.Filters.MinScore == 0 { cfg.Filters.MinScore = 0 }
	return cfg, nil
}
