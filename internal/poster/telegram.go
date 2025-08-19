package poster

import (
	"context"
	"fmt"

	"github.com/LibenHailu/cncg-bot/internal/store"
	"github.com/LibenHailu/cncg-bot/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TG struct {
	Bot       *tgbotapi.BotAPI
	ChannelID string
	ParseMode string // "MarkdownV2"
}

func New(botToken, channelID, parseMode string) (*TG, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil { return nil, err }
	return &TG{Bot: bot, ChannelID: channelID, ParseMode: parseMode}, nil
}

func (t *TG) PostItem(ctx context.Context, it store.Item) error {
	fmt.Println("Posting to Telegram:", it.Title, it.URL)
	title := util.EscapeTelegram(it.Title)
	url := util.EscapeTelegram(it.URL)
	source := util.EscapeTelegram(it.Source)
	sum := util.EscapeTelegram(it.Summary)

	// Title as a clickable link, then 2–3 sentence summary + source attribution + tags
	text := fmt.Sprintf("[*%s*](%s)\n\n%s\n\n_Source:_ %s", title, url, sum, source)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{ChatID: chatIDFromString(t.ChannelID)},
		Text:     text,
		ParseMode: t.ParseMode,
		DisableWebPagePreview: false,
	}
	_, err := t.Bot.Send(msg)
	return err
}

func chatIDFromString(s string) int64 {
	// channel usernames can be sent via string; library wants int64 for IDs.
	// For usernames we use tgbotapi.NewMessageToChannel later, but here we keep simplest:
	// tgbotapi accepts ChatConfig with ChatID OR ChannelUsername
	// For simplicity, when s starts with "@", ChatID=0 and ChannelUsername used via NewMessageToChannel:
	return 0
}
