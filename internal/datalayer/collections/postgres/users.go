package postgres

import (
	"context"
	"fmt"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

type UsersCollection struct {
	users map[int64]models.Users
}

func (c *UsersCollection) Create(ctx context.Context, userRequest models.Users) error {
	_, ok := c.users[userRequest.ChatID]
	if ok {
		return nil
	}

	c.users[userRequest.ChatID] = userRequest

	fmt.Println("user created: ", userRequest)

	return nil
}
func (ac *models.Users) SetActive(ctx context.Context, active bool) {

}
func (c *UsersCollection) GetByChatID(ctx context.Context, chatID int64) (*models.Users, error) {
	user, ok := c.users[chatID]
	if ok {
		return &user, nil
	}

	return nil, nil
}
