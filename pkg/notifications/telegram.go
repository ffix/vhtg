package notifications

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	tokenEmpty  = errors.New("TELEGRAM_BOT_TOKEN environment variable must be set")
	chatIDEmpty = errors.New("TELEGRAM_CHAT_ID environment variable must be set")
)

type TelegramClient struct {
	chatID int64
	api    *tgbotapi.BotAPI
	userID int64
	logger logger
}

func NewTelegramClient(logger logger) (*TelegramClient, error) {
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

	user, err := api.GetMe()
	if err != nil {
		return nil, fmt.Errorf("failed to get my own user id: %w", err)
	}

	return &TelegramClient{chatID: chatID, api: api, userID: user.ID, logger: logger}, nil
}

func (t *TelegramClient) SendMessage(msg string, pin bool) error {
	var err error

	chatConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: t.chatID}}

	for {
		if !pin {

			break
		}
		t.logger.Debug("Getting chat info...")
		chat, err := t.api.GetChat(chatConfig)
		if err != nil {
			return fmt.Errorf("faield to get chat info: %w", err)
		}

		pinned := chat.PinnedMessage
		t.logger.Debugf("got pinned message: %v", pinned)
		if pinned == nil || pinned.From.ID != t.userID {
			break
		}

		t.logger.Debugf("Unpinning message %d...", pinned.MessageID)
		msgConfig := tgbotapi.UnpinChatMessageConfig{MessageID: pinned.MessageID, ChatID: t.chatID}
		if _, err := t.api.Request(msgConfig); err != nil {
			return fmt.Errorf("failed to unpin message :%w", err)
		}
		// We need this delay because Telegram doesn't return a new recent pinned message without it.
		time.Sleep(5 * time.Second)
	}

	msgConfig := tgbotapi.NewMessage(t.chatID, msg)

	t.logger.Debug("Senging message...")
	message, err := t.api.Send(msgConfig)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	t.logger.Debugf("Sent a message with id %d", message.MessageID)

	if !pin {
		return nil
	}

	t.logger.Debug("Pinning message...")
	pinConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              t.chatID,
		MessageID:           message.MessageID,
		DisableNotification: true,
	}

	if _, err = t.api.Request(pinConfig); err != nil {
		return fmt.Errorf("failed to pin message via Telegram Bot API: %w", err)
	}
	return nil
}
