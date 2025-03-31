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
	if UserRequest.Backload <= 0 {
		return fmt.Sprintf(
			"Результаты расчета:\n"+
				"Пройденное расстояние за день: %d км\n"+
				"Общий расход топлива: %g л\n"+
				"Остаток топлива на конец дня: %g л\n"+
				"Пробег на конец дня: %d км\n"+
				`=================
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
	} else {
		return fmt.Sprintf(
			"Результаты расчета:\n"+
				"Пройденное расстояние за день: %d км\n"+
				"Общий расход топлива: %g л\n"+
				"Остаток топлива на конец дня: %g л\n"+
				"Пробег на конец дня: %d км\n"+
				`=================
				Итог:
				Расход на пробег:  %v*%g/100=%g
				Pасход на подъемы: %v*%g=%g
				Pасход c обратным: +%g*%v*%v/100=%g
				=================`,

			v.OperatingDistance,
			v.TotalFuel,
			v.DailyRate,
			v.DailyRun,
			v.OperatingDistance,
			UserRequest.Consumption,
			v.Wastage,
			(UserRequest.QuantityTrips)*2,
			UserRequest.Lifting,
			v.Lifting,
			v.Undelivery,
			UserRequest.QuantityTrips,
			UserRequest.Distance,
			v.Underfuel,
		)
	}

}
