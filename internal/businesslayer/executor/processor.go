package executor

import (
	"math"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
)

type ExecutorUSerRequest struct {
	dto.UserRequest
}

func (u ExecutorUSerRequest) Executor() dto.VitalData {
	var v dto.VitalData

	v.Undelivery = float64(u.Capacity - u.Tons)                                                      //недотонны
	v.OperatingDistance = u.Distance * 2 * u.QuantityTrips                                           // Пройденное расстояние за день
	v.Wastage = roundTo(float64(v.OperatingDistance)*float64(u.Consumption/100), 1)                  // Расход топлива на эти километры
	v.Lifting = float64(u.QuantityTrips) * u.Lifting                                                 // Подъемы
	v.Underfuel = roundTo(float64(v.Undelivery)*float64(u.QuantityTrips)*float64(u.Distance)/100, 1) // Расход топлива на недовоз
	v.TotalFuel = roundTo(v.Wastage+v.Lifting-v.Underfuel, 1)                                        // Общий расход топлива
	v.DailyRate = roundTo(float64(u.FuelResidue)+float64(u.Refuel)-v.TotalFuel, 1)                   // Расход на день с учетом заправки
	v.DailyRun = u.SpeedometerResidue + v.OperatingDistance                                          // Пробег за день

	if u.Backload > 0 {
		v.Lifting = 2 * float64(u.QuantityTrips) * u.Lifting                                             // Подъемы
		v.Undelivery = math.Max(0, float64(u.Tons+u.Backload-u.Capacity))                                // перевоз тонн
		v.Underfuel = roundTo(float64(v.Undelivery)*float64(u.QuantityTrips)*float64(u.Distance)/100, 1) // расход топлива за день
		v.TotalFuel = roundTo(v.Wastage+v.Lifting+v.Underfuel, 1)                                        // общий расход топлива
		v.DailyRate = roundTo(float64(u.FuelResidue)+float64(u.Refuel)-v.TotalFuel, 1)                   // расход на день с учетом заправки
	}

	return v
}

func roundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
