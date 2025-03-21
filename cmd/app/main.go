package main

import (
	"context"
	"os"

	"Project-Chat_Bot-for-drivergithub.com/PapaDjo2000/Project-Chat_Bot-for-drivers/interanl/datalayer/collections/cache"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/interanl/businesslayer/domain/bot"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/interanl/businesslayer/domain/users"

	"github.com/rs/zerolog"
)

// tgBot - @BotFather
// @GetMyChatID_BestBot

// THIS VALUES SHOULD BE IN CONFIG/ENV FILE
const (
	apiKey      = "8008625848:AAFQ-xFfdNo7KS3cBM0JBqbNvS-bDEmDkzI"
	adminChatID = 401631302
)

func main() {
	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	usersCollection := cache.NewUsersCollection()
	usersProcessor := users.NewProcessor(logger, usersCollection)

	tgBot, err := bot.New(apiKey, logger, usersProcessor)
	if err != nil {
		logger.Err(err).Send()

		return
	}

	if err := tgBot.SendMessage(adminChatID, "bot started"); err != nil {
		logger.Err(err).Send()

		return
	}

	defer func() {
		if err := tgBot.SendMessage(adminChatID, "bot finished"); err != nil {
			logger.Err(err).Send()

			return
		}
	}()

	if err := tgBot.Listen(ctx); err != nil {
		logger.Err(err).Send()

		return
	}
}
