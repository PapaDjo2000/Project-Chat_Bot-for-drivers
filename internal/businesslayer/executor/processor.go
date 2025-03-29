package executor

import (
	"math"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Calculate(UserRequest dto.UserRequest) dto.VitalData {
	var v dto.VitalData
	v.Undelivery = float64(UserRequest.Capacity - UserRequest.Tons)                                                      //недотонны
	v.OperatingDistance = UserRequest.Distance * 2 * UserRequest.QuantityTrips                                           // Пройденное расстояние за день
	v.Wastage = roundTo(float64(v.OperatingDistance)*float64(UserRequest.Consumption/100), 1)                            // Расход топлива на эти километры
	v.Lifting = float64(UserRequest.QuantityTrips) * UserRequest.Lifting                                                 // Подъемы
	v.Underfuel = roundTo(float64(v.Undelivery)*float64(UserRequest.QuantityTrips)*float64(UserRequest.Distance)/100, 1) // Расход топлива на недовоз
	v.TotalFuel = roundTo(v.Wastage+v.Lifting-v.Underfuel, 1)                                                            // Общий расход топлива
	v.DailyRate = roundTo(float64(UserRequest.FuelResidue)+float64(UserRequest.Refuel)-v.TotalFuel, 1)                   // Расход на день с учетом заправки
	v.DailyRun = UserRequest.SpeedometerResidue + v.OperatingDistance                                                    // Пробег за день

	if UserRequest.Backload > 0 {
		v.Lifting = 2 * float64(UserRequest.QuantityTrips) * UserRequest.Lifting                                             // Подъемы
		v.Undelivery = math.Max(0, float64(UserRequest.Tons+UserRequest.Backload-UserRequest.Capacity))                      // перевоз тонн
		v.Underfuel = roundTo(float64(v.Undelivery)*float64(UserRequest.QuantityTrips)*float64(UserRequest.Distance)/100, 1) // расход топлива за день
		v.TotalFuel = roundTo(v.Wastage+v.Lifting+v.Underfuel, 1)                                                            // общий расход топлива
		v.DailyRate = roundTo(float64(UserRequest.FuelResidue)+float64(UserRequest.Refuel)-v.TotalFuel, 1)                   // расход на день с учетом заправки
	}
	return v
}
func roundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
