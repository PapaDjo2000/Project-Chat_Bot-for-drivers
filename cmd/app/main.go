package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/domain/bot"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/domain/users"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections/postgres"

	"github.com/rs/zerolog"
)

// tgBot - @BotFather
// @GetMyChatID_BestBot

// THIS VALUES SHOULD BE IN CONFIG/ENV FILE

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
		return
	}
	adminChatIDStr := os.Getenv("ADMIN_CHAT_ID")
	if adminChatIDStr == "" {
		log.Fatalf("ADMIN_CHAT_ID is not set")
		return
	}
	adminChatID, err := strconv.ParseInt(adminChatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse ADMIN_CHAT_ID: %v\n", err)
		return
	}
	if err != nil {
		log.Fatalf("Failed to parse ADMIN_CHAT_ID: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	usersCollection := postgres.NewUserStorage(db)
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
