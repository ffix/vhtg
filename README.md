# VHTG - Valheim to Telegram Message Processor

VHTG (Valheim to Telegram Message Processor) is a command-line application written in Go that reads a text message from the standard input, matches it using a regular expression, and sends a formatted message to a specified Telegram chat using the Telegram Bot API. The primary use case for this application is to process and forward messages from a Valheim server. The application also supports a test mode that only outputs the formatted message to the standard output without sending it to Telegram.

## Prerequisites

- Go 1.18 or higher
- A Telegram bot token (you can create a new bot by talking to the [BotFather](https://core.telegram.org/bots#6-botfather))
- A Telegram chat ID to send messages to (you can use a personal chat, group chat, or channel)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/vhtg.git
cd vhtg
```

2. Build the application:

```bash
go build -o vhtg main.go
```

## Usage

1. Set the `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` environment variables:

```bash
export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
export TELEGRAM_CHAT_ID="your-telegram-chat-id"
```

2. Run the application, providing the regular expression and output pattern as command-line arguments:

```bash
echo 'Your Valheim server log message here' | ./vhtg <regex_pattern> <output_pattern>
```

For example:

```bash
echo 'Session "MyServer" with join code 547281 and IP 127.0.0.1:2456 is active with 0 player(s)' | ./vhtg 'Session "(?P<session_name>\w+)" with join code (?P<join_code>\d+)' 'Session "{session_name}" has join code {join_code}'
```

3. To enable test mode and only output the formatted message to the standard output without sending it to Telegram, add the `--test` flag:

```bash
echo 'Your Valheim server log message here' | ./vhtg --test <regex_pattern> <output_pattern>
```

4. To suppress error messages when no matches are found, add the --quiet flag:

```bash
echo 'Your Valheim server log message here' | ./vhtg --quiet <regex_pattern> <output_pattern>
```

## License

This project is licensed under the MIT License.
