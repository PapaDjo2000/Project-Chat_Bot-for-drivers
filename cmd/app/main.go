package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/domain/bot"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/domain/users"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections"

	"github.com/rs/zerolog"
)

// tgBot - @BotFather
// @GetMyChatID_BestBot

// THIS VALUES SHOULD BE IN CONFIG/ENV FILE
const (
	apiKey      = "7686022156:AAHpISLuOFWkUksQBLcUGfx0dEiggcBs-OA"
	adminChatID = 524060834
)

func main() {

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	usersCollection := collections.Users.CreateUser(ctx)
	usersProcessor := users.NewProcessor(logger, usersCollection)

	tgBot, err := bot.New(token, logger, usersProcessor)
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
