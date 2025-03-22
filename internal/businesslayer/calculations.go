package businesslayer

import (
	"fmt"
	"math"

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

func (ur *UserRequest) SetLifting(lif int) (int, error) {
	if lif < 0 {
		return 0, fmt.Errorf("invalid consumption value")
	}
	ur.Backload = lif
	return ur.Backload, nil
}

func (ur *UserRequest) Сalculations() VitalData {
	var vd VitalData

	vd.Undelivery = float64(ur.Capacity - ur.Tons)                                                       //недотонны
	vd.OperatingDistance = ur.Distance * 2 * ur.QuantityTrips                                            // Пройденное расстояние за день
	vd.Wastage = roundTo(float64(vd.OperatingDistance)*float64(ur.Consumption/100), 1)                   // Расход топлива на эти километры
	vd.Lifting = float64(ur.QuantityTrips) * ur.Lifting                                                  // Подъемы
	vd.Underfuel = roundTo(float64(vd.Undelivery)*float64(ur.QuantityTrips)*float64(ur.Distance)/100, 1) // Расход топлива на недовоз
	vd.TotalFuel = roundTo(vd.Wastage+vd.Lifting-vd.Underfuel, 1)                                        // Общий расход топлива
	vd.DailyRate = roundTo(float64(ur.FuelResidue)+float64(ur.Refuel)-vd.TotalFuel, 1)                   // Расход на день с учетом заправки
	vd.DailyRun = ur.SpeedometerResidue + vd.OperatingDistance                                           // Пробег за день

	if ur.Backload > 0 {
		vd.Lifting = 2 * float64(ur.QuantityTrips) * 0.5                                                     // Подъемы
		vd.Undelivery = math.Max(0, float64(ur.Tons+ur.Backload-ur.Capacity))                                // перевоз тонн
		vd.Underfuel = roundTo(float64(vd.Undelivery)*float64(ur.QuantityTrips)*float64(ur.Distance)/100, 1) // расход топлива за день
		vd.TotalFuel = roundTo(vd.Wastage+vd.Lifting+vd.Underfuel, 1)                                        // общий расход топлива
		vd.DailyRate = roundTo(float64(ur.FuelResidue)+float64(ur.Refuel)-vd.TotalFuel, 1)                   // расход на день с учетом заправки
	}

	return vd
}

func roundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
