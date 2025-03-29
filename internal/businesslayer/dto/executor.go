package dto

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRequest struct {
	Consumption        float64 // расход
	Capacity           int     // тонны грузоподьемности
	FuelResidue        float64 // остаток топлива
	SpeedometerResidue int     // Остаток по спидометру
	Refuel             int     // заправка
	Distance           int     // растояние в одну сторону
	QuantityTrips      int     // количество рейсов
	Tons               int     // тоны
	Backload           int     // обратные тонны
	Lifting            float64 // подьемы
}

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

func NewVitalData() *VitalData {
	return &VitalData{}
}

func (v *VitalData) ToString(UserRequest UserRequest) string {
	return fmt.Sprintf(
		"Результаты расчета:\n"+
			"Недотонны: %.1f\n"+
			"Пройденное расстояние за день: %d км\n"+
			"Расход топлива: %.1f л\n"+
			"Подъемы: %.1f л\n"+
			"Расход топлива на недовоз: %.1f л\n"+
			"Общий расход топлива: %.1f л\n"+
			"Остаток топлива на конец дня: %.1f л\n"+
			"Пробег на конец дня: %d км\n"+
			`=================
			%v*%2.f=%1.f
			%v*%1.f=%1.f
			%1.f*%v*%v/100=%1.f
			=================`,
		v.Undelivery,
		v.OperatingDistance, // остаётся int
		v.Wastage,
		v.Lifting,
		v.Underfuel,
		v.TotalFuel,
		v.DailyRate,
		v.DailyRun, // остаётся int
		v.OperatingDistance,
		UserRequest.Consumption,
		v.Wastage,
		UserRequest.QuantityTrips,
		UserRequest.Lifting,
		v.Lifting,
		v.Undelivery,
		UserRequest.QuantityTrips,
		UserRequest.Distance,
		v.Underfuel,
	)

}
