package models

import (
	"time"

	"github.com/google/uuid"
)

type Reports struct {
	Id       uuid.UUID
	User_id  uuid.UUID
	Date     time.Time
	Request  string
	Response string
}
