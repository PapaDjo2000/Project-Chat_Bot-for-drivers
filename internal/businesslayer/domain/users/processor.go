package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Processor struct {
	logger          zerolog.Logger
	usersCollection collections.Users
}

func NewProcessor(logger zerolog.Logger, usersCollection collections.Users) *Processor {
	return &Processor{logger: logger, usersCollection: usersCollection}
}

func (p *Processor) CreateIfNotExist(ctx context.Context, userRequest dto.User) error {
	user, err := p.usersCollection.GetUserByChatID(ctx, userRequest.ChatID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			p.logger.Err(err).Msg("Database error while fetching user")
			return fmt.Errorf("failed to fetch user: %w", err)
		}
	}

	if user == nil {
		if err := p.usersCollection.CreateUser(
			ctx,
			&models.Users{
				ID:     uuid.New(),
				Name:   userRequest.Name,
				ChatID: userRequest.ChatID,
			},
		); err != nil {
			p.logger.Err(err).Send()
			return err
		}
	}

	return nil
}

func (p *Processor) LoadByChatID(ctx context.Context, chatID int64) (*dto.User, error) {
	user, err := p.usersCollection.GetUserByChatID(ctx, chatID)
	if err != nil {
		p.logger.Err(err).Send()
		return nil, err
	}

	if user == nil {
		return nil, errors.New("not found")
	}

	return &dto.User{
		ID:     user.ID,
		Name:   user.Name,
		ChatID: user.ChatID,
	}, nil
}
