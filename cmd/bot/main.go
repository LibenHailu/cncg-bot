package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LibenHailu/cncg-bot/internal/config"
	"github.com/LibenHailu/cncg-bot/internal/core"
	"github.com/LibenHailu/cncg-bot/internal/poster"
	"github.com/LibenHailu/cncg-bot/internal/schedule"
	"github.com/LibenHailu/cncg-bot/internal/store"
)

func main() {
	cfgPath := os.Getenv("BOT_CONFIG")
	if cfgPath == "" { cfgPath = "config.yaml" }

	cfg, err := config.Load(cfgPath)
	if err != nil { log.Fatal("config:", err) }

	db, err := store.Open(cfg.DBPath)
	if err != nil { log.Fatal("db:", err) }

	p := core.Pipeline{
		Filters: core.Filters{
			MaxAgeDays: cfg.Filters.MaxAgeDays,
			MinScore:   cfg.Filters.MinScore,
			Positive:   cfg.Keywords.Positive,
			Negative:   cfg.Keywords.Negative,
		},
		DB: db,
	}
	for _, s := range cfg.Sources {
		p.Sources = append(p.Sources, core.SourceCfg{
			Name: s.Name, Type: s.Type, URL: s.URL, Weight: s.Weight, Tags: s.Tags,
		})
	}

	tg, err := poster.New(cfg.Telegram.BotToken, cfg.Telegram.ChannelID, cfg.Telegram.ParseMode)
	if err != nil { log.Fatal("telegram:", err) }

	// warm run on startup
	ctx := context.Background()
	if err := p.RunOnce(ctx); err != nil {
		db.LogError(ctx,"pipeline", err.Error())
	}

	// start scheduler
	job := schedule.Job{
		CronSpec: cfg.Scheduler.CronSpec,
		BatchSize: cfg.Scheduler.BatchSize,
		DB: db,
		Poster: tg,
		MinScore: cfg.Filters.MinScore,
	}
	c, err := job.Start(ctx)
	if err != nil { log.Fatal("cron:", err) }
	log.Println("bot started")

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down...")
	c.Stop()
	time.Sleep(500 * time.Millisecond)
}
