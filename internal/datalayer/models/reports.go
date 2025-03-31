package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Reports struct {
	ID       uuid.UUID
	UserID   int64
	Date     time.Time
	Request  json.RawMessage
	Response json.RawMessage
}
