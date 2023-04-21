package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/pflag"
)

const maxRetries = 3
const retryDelay = 5 * time.Second

func sendTelegramMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		_, err = bot.Send(msg)
		if err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to send message via Telegram Bot API after %d retries: %v", maxRetries, err)
}

func processMessage(inputText, regexPattern, outputPattern string, testMode bool, quietMode bool) error {
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindStringSubmatch(inputText)

	if len(matches) == 0 {
		if quietMode {
			return nil
		}
		return errors.New("no matching message found")
	}

	// Replace the placeholders in the output pattern with the matched values
	outputMessage := outputPattern
	for i, name := range regex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		outputMessage = regexp.MustCompile("{"+name+"}").ReplaceAllString(outputMessage, matches[i])
	}

	if testMode {
		// Output the message to stdout without sending it to Telegram
		fmt.Println(outputMessage)
		return nil
	}

	// Set up the Telegram bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatIDStr == "" {
		return errors.New("TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables must be set")
	}

	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create Telegram bot: %v", err)
	}

	// Send the message with retries
	msg := tgbotapi.NewMessage(chatID, outputMessage)
	err = sendTelegramMessage(bot, msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var testMode, quietMode bool

	pflag.BoolVar(&testMode, "test", false, "Enable test mode (don't send messages to Telegram)")
	pflag.BoolVar(&quietMode, "quiet", false, "Enable quiet mode (don't report an error if there is no match)")
	pflag.Parse()

	args := pflag.Args()

	if len(args) != 2 {
		fmt.Println("Usage: go run main.go [--test] [--quiet] <regex_pattern> <output_pattern>")
		os.Exit(1)
	}

	regexPattern := args[0]
	outputPattern := args[1]

	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	err := processMessage(inputText, regexPattern, outputPattern, testMode, quietMode)
	if err != nil {
		fmt.Printf("Error processing message: %v\n", err)
		os.Exit(1)
	}
}
