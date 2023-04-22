package notifications

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	tokenEmpty  = errors.New("TELEGRAM_BOT_TOKEN environment variable must be set")
	chatIDEmpty = errors.New("TELEGRAM_CHAT_ID environment variable must be set")
)

const maxRetries = 3
const retryDelay = 5 * time.Second

type TelegramClient struct {
	chatID int64
	api    *tgbotapi.BotAPI
}

func NewTelegramClient() (*TelegramClient, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" {
		return nil, tokenEmpty
	}
	if chatIDStr == "" {
		return nil, chatIDEmpty
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid Chat ID: %w", err)
	}

	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot api: %w", err)
	}

	return &TelegramClient{chatID: chatID, api: api}, nil
}

func (t *TelegramClient) SendMessage(msg string) error {
	var err error

	msgConfig := tgbotapi.NewMessage(t.chatID, msg)

	for i := 0; i < maxRetries; i++ {
		_, err = t.api.Send(msgConfig)
		if err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send message via Telegram Bot API after %d retries: %v", maxRetries, err)
}
