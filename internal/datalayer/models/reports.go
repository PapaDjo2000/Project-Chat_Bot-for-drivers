package models

import (
	"time"

	"github.com/google/uuid"
)

type Reports struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Date     time.Time
	Request  string
	Response string
}
