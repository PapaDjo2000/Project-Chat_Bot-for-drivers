package models

import (
	"github.com/google/uuid"
)

type Users struct {
	Id      uuid.UUID
	Name    string
	Chat_id int
	Active  bool
}
