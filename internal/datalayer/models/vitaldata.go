package models

import "github.com/google/uuid"

type VitalData struct {
	UserId            uuid.UUID
	Undelivery        float64
	OperatingDistance int
	Wastage           float64
	Lifting           float64
	Underfuel         float64
	TotalFuel         float64
	DailyRun          int
	DailyRate         float64
}
