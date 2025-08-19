package schedule

import (
	"context"
	"log"

	"github.com/LibenHailu/cncg-bot/internal/poster"
	"github.com/LibenHailu/cncg-bot/internal/store"
	"github.com/robfig/cron/v3"
)

type Job struct {
	CronSpec   string
	BatchSize  int
	DB         *store.Store
	Poster     *poster.TG
	MinScore   float64
}

func (j *Job) Start(ctx context.Context) (*cron.Cron, error) {
	log.Println("Starting scheduler with cron spec:", j.CronSpec)
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 "+j.CronSpec, func() {
		j.runOnce(ctx)
	})
	if err != nil { return nil, err }
	c.Start()
	return c, nil
}

func (j *Job) runOnce(ctx context.Context) {
	items, err := j.DB.NextUnposted(ctx, j.MinScore, j.BatchSize)
	if err != nil { j.DB.LogError(ctx,"schedule:select", err.Error()); return }
	for _, it := range items {
		log.Println("Posting item to Telegram:", it.Title, it.URL)
		if err := j.Poster.PostItem(ctx, it); err != nil {
			j.DB.LogError(ctx,"telegram:send", err.Error())
			continue
		}
		if err := j.DB.MarkPosted(ctx, it.ID); err != nil {
			log.Println("mark posted error:", err)
		}
	}
}
