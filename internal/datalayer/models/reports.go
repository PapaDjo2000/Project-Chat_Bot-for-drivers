package models

import (
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/google/uuid"
)

type Reports struct {
	ID       uuid.UUID
	UserID   int64
	Date     time.Time
	Request  dto.UserRequest
	Response string
}
