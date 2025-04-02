package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVitalData_ToString(t *testing.T) {
	tests := []struct {
		name           string
		vitalData      *VitalData
		userRequest    UserRequest
		expectedOutput string
	}{
		{
			name: "1",
			vitalData: &VitalData{
				UserId:            uuid.New(),
				Undelivery:        10.0,
				OperatingDistance: 200,
				Wastage:           16.0,
				Lifting:           3.0,
				Underfuel:         20.0,
				TotalFuel:         60.0,
				DailyRun:          300,
				DailyRate:         10.0,
			},
			userRequest: UserRequest{
				Consumption:        8.0,
				Capacity:           10,
				FuelResidue:        20.0,
				SpeedometerResidue: 500,
				Refuel:             40,
				Distance:           100,
				QuantityTrips:      2,
				Tons:               5,
				Backload:           0,
				Lifting:            1.5,
			},
			expectedOutput: `Результаты расчета:
			Пройденное расстояние за день: 200 км
			Общий расход топлива: 60 л
			Остаток топлива на конец дня: 10 л
			Пробег на конец дня: 300 км
			=================
			Итог:
			Расход на пробег:   200*8/100=16
			Pасход на подъемы:  2*1.5=3
			Pасход с недовоз : -10*2*100/100=20
			=================`,
		},
		{
			name: "2",
			vitalData: &VitalData{
				UserId:            uuid.New(),
				Undelivery:        5.0,
				OperatingDistance: 150,
				Wastage:           15.0,
				Lifting:           12.0,
				Underfuel:         12.0,
				TotalFuel:         50.0,
				DailyRun:          250,
				DailyRate:         15.0,
			},
			userRequest: UserRequest{
				Consumption:        10.0,
				Capacity:           12,
				FuelResidue:        25.0,
				SpeedometerResidue: 400,
				Refuel:             30,
				Distance:           80,
				QuantityTrips:      3,
				Tons:               6,
				Backload:           2,
				Lifting:            2.0,
			},
			expectedOutput: `Результаты расчета:
			Пройденное расстояние за день: 150 км
			Общий расход топлива: 50 л
			Остаток топлива на конец дня: 15 л
			Пробег на конец дня: 250 км
			=================
			Итог:
			Расход на пробег:   150*10/100=15
			Pасход на подъемы:  6*2=12
			Pасход c обратным: +5*3*80/100=12
			=================`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vitalData.ToString(tt.userRequest)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
