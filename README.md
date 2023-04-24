# VHTG - Valheim to Telegram Message Processor

VHTG is a command-line application written in Go that reads text messages from the standard input or a Docker container with the Valheim server, matches them using regular expressions, and sends formatted messages to a specified Telegram chat using the Telegram Bot API.

## Features

- Automatically process Valheim server logs
- Support for reading logs from standard input or a Docker container
- Send formatted messages to a specified Telegram chat
- Test mode for outputting formatted messages without sending them to Telegram

## Prerequisites

- Go 1.18 or higher
- A Telegram bot token (you can create a new bot by talking to the [BotFather](https://core.telegram.org/bots#6-botfather))
- A Telegram chat ID to send messages to (you can use a personal chat, group chat, or channel)

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/ffix/vhtg.git
cd vhtg
```

## Build

Use the provided Makefile to build the application for your platform:

```bash
make vhtg # For your current platform
make all  # For all supported platforms
```

## Usage

1. Set the `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, and `SERVER_PASS` environment variables:

```bash
export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
export TELEGRAM_CHAT_ID="your-telegram-chat-id"
export SERVER_PASS="your-server-password"
```

2. In standard mode, pipe the application from the Valheim server:

```bash
./valheim-server | ./vhtg
```

3. To enable test mode and only output the formatted message to the standard output without sending it to Telegram, add the `--test` flag:

```bash
./valheim-server | ./vhtg --test
```

4. To integrate with Docker and poll the logs internally, add the `--docker` flag:

```bash
./vhtg --docker
```

### Configuration

The VHTG application is primarily configured using environment variables. These include:

- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token
- `TELEGRAM_CHAT_ID`: The Telegram chat ID to send messages to
- `SERVER_PASS`: Your Valheim server password

### Running with Docker Compose

To run the application inside a Docker container, you can use the following Docker Compose configuration:

```yaml
services:
  vhtg:
    build:
      context: vhtg
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    env_file:
      - "env.valheim"
    restart: always
    stop_grace_period: 2m
```

## Contributing

Contributions to the VHTG project are welcome! If you'd like to contribute, please follow these steps:

1. Fork the repository
2. Create a new branch for your changes (`git checkout -b my-feature`)
3. Commit your changes (`git commit -am 'Add my feature'`)
4. Push the changes to the branch (`git push origin my-feature`)
5. Create a new pull request

Please make sure to follow the project's coding standards and include tests for any new features or bug fixes.

## Support

If you encounter any issues or need help using the VHTG application, please:

1. Check the project's documentation and README for any relevant information
2. Search the project's issue tracker for any similar issues
3. If you cannot find a solution, create a new issue with a clear description of the problem, steps to reproduce it, and any relevant logs or error messages

We will do our best to address your concerns and provide assistance.

## License

This project is licensed under the MIT License.
