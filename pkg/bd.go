package main

import (
	"github.com/google/uuid"
)

type Users struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
}
type Requests struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	OrderData string    `json:"order_date"`
	Data      string    `json:"data"`
}
