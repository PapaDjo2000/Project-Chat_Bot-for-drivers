package executor

import (
	"math"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Calculate(userRequest dto.UserRequest) dto.VitalData {
	var v dto.VitalData
	v.Undelivery = float64(userRequest.Capacity - userRequest.Tons)
	v.OperatingDistance = userRequest.Distance * 2 * userRequest.QuantityTrips
	v.Wastage = roundTo(float64(v.OperatingDistance) * float64(userRequest.Consumption/100))
	v.Lifting = float64(userRequest.QuantityTrips) * userRequest.Lifting
	v.Underfuel = roundTo(float64(v.Undelivery) * float64(userRequest.QuantityTrips) * float64(userRequest.Distance) / 100)
	v.TotalFuel = roundTo(v.Wastage + v.Lifting - v.Underfuel)
	v.DailyRate = roundTo(float64(userRequest.FuelResidue) + float64(userRequest.Refuel) - v.TotalFuel)
	v.DailyRun = userRequest.SpeedometerResidue + v.OperatingDistance

	if userRequest.Backload > 0 {
		v.Lifting = 2 * float64(userRequest.QuantityTrips) * userRequest.Lifting
		v.Undelivery = math.Max(0, float64(userRequest.Tons+userRequest.Backload-userRequest.Capacity))
		v.Underfuel = roundTo(float64(v.Undelivery) * float64(userRequest.QuantityTrips) *
			float64(userRequest.Distance) / 100)
		v.TotalFuel = roundTo(v.Wastage + v.Lifting + v.Underfuel)
		v.DailyRate = roundTo(float64(userRequest.FuelResidue) + float64(userRequest.Refuel) - v.TotalFuel)
	}
	return v
}

func roundTo(value float64) float64 {
	return math.Round(value*10) / 10
}
