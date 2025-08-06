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
	UserID            uuid.UUID
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

func (v *VitalData) ToString(userRequest UserRequest) string {
	if userRequest.Backload <= 0 {
		return fmt.Sprintf(
			`Результаты расчета:
			Пройденное расстояние за день: %d км
			Общий расход топлива: %g л
			Остаток топлива на конец дня: %g л
			Пробег на конец дня: %d км
			=================
			Итог:
			Расход на пробег:   %v*%g/100=%g
			Pасход на подъемы:  %v*%g=%g
			Pасход с недовоз : -%g*%v*%v/100=%g
			=================`,

			v.OperatingDistance,
			v.TotalFuel,
			v.DailyRate,
			v.DailyRun,
			v.OperatingDistance,
			userRequest.Consumption,
			v.Wastage,
			userRequest.QuantityTrips,
			userRequest.Lifting,
			v.Lifting,
			v.Undelivery,
			userRequest.QuantityTrips,
			userRequest.Distance,
			v.Underfuel,
		)
	}
	return fmt.Sprintf(
		`Результаты расчета:
			Пройденное расстояние за день: %d км
			Общий расход топлива: %g л
			Остаток топлива на конец дня: %g л
			Пробег на конец дня: %d км
			=================
			Итог:
			Расход на пробег:   %v*%g/100=%g
			Pасход на подъемы:  %v*%g=%g
			Pасход c обратным: +%g*%v*%v/100=%g
			=================`,

		v.OperatingDistance,
		v.TotalFuel,
		v.DailyRate,
		v.DailyRun,
		v.OperatingDistance,
		userRequest.Consumption,
		v.Wastage,
		(userRequest.QuantityTrips)*2,
		userRequest.Lifting,
		v.Lifting,
		v.Undelivery,
		userRequest.QuantityTrips,
		userRequest.Distance,
		v.Underfuel,
	)
}
