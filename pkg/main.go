package main

import (
	"fmt"
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

func (ur *UserRequest) SetConsumption(cons float64) (float64, error) {
	if cons <= 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Consumption = cons
	return ur.Consumption, nil
}
func (ur *UserRequest) SetCapacity(cap int) (int, error) {
	if cap <= 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Capacity = cap
	return ur.Capacity, nil
}
func (ur *UserRequest) SetFuelResidue(fr float64) (float64, error) {
	if fr < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.FuelResidue = fr
	return ur.FuelResidue, nil
}
func (ur *UserRequest) SetSpeedometerResidue(spr int) (int, error) {
	if spr < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.SpeedometerResidue = spr
	return ur.SpeedometerResidue, nil
}

func (ur *UserRequest) SetRefuel(ref int) (int, error) {
	if ref < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Refuel = ref
	return ur.Refuel, nil
}

func (ur *UserRequest) SetDistance(dis int) (int, error) {
	if dis < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Distance = dis
	return ur.Distance, nil
}
func (ur *UserRequest) SetQuantityTrips(qt int) (int, error) {
	if qt <= 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.QuantityTrips = qt
	return ur.QuantityTrips, nil
}
func (ur *UserRequest) SetTons(ton int) (int, error) {
	if ton <= 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Tons = ton
	return ur.Tons, nil
}
func (ur *UserRequest) SetBackload(dton int) (int, error) {
	if dton < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Backload = dton
	return ur.Backload, nil
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
