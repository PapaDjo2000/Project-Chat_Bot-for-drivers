package main

import (
	"math"
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
}

func NewUserRequest(Consumption, FuelResidue float64, Capacity, SpeedometerResidue, Refuel, Distance, QuantityTrips, Tons, Backload int) *UserRequest {
	return &UserRequest{
		Consumption:        Consumption,
		Capacity:           Capacity,
		FuelResidue:        FuelResidue,
		SpeedometerResidue: SpeedometerResidue,
		Refuel:             Refuel,
		Distance:           Distance,
		QuantityTrips:      QuantityTrips,
		Tons:               Tons,
		Backload:           Backload,
	}
}

func (ur *UserRequest) Сalculations() (int, float64) {
	undelivery := float64(ur.Capacity - ur.Tons)                                                    //недотонны
	operatingDistance := ur.Distance * 2 * ur.QuantityTrips                                         // Пройденное расстояние за день
	wastage := roundTo(float64(operatingDistance)*float64(ur.Consumption/100), 1)                   // Расход топлива на эти километры
	lifting := float64(ur.QuantityTrips) * 0.5                                                      // Подъемы
	underfuel := roundTo(float64(undelivery)*float64(ur.QuantityTrips)*float64(ur.Distance)/100, 1) // Расход топлива на недовоз
	totalFuel := roundTo(wastage+lifting-underfuel, 1)                                              // Общий расход топлива
	dailyRate := roundTo(float64(ur.FuelResidue)+float64(ur.Refuel)-totalFuel, 1)                   // Расход на день с учетом заправки
	dailyRun := ur.SpeedometerResidue + operatingDistance                                           // Пробег за день

	if ur.Backload > 0 {
		lifting = 2 * float64(ur.QuantityTrips) * 0.5                                                  // Подъемы
		undelivery = math.Max(0, float64(ur.Tons+ur.Backload-ur.Capacity))                             // перевоз тонн
		underfuel = roundTo(float64(undelivery)*float64(ur.QuantityTrips)*float64(ur.Distance)/100, 1) // расход топлива за день
		totalFuel = roundTo(wastage+lifting+underfuel, 1)                                              // общий расход топлива
		dailyRate = roundTo(float64(ur.FuelResidue)+float64(ur.Refuel)-totalFuel, 1)                   // расход на день с учетом заправки
	}

	return dailyRun, dailyRate
}

func roundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
