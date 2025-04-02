package executor

import (
	"testing"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/stretchr/testify/assert"
)

func TestProcessor_Calculate(t *testing.T) {
	tests := []struct {
		name         string
		userRequest  dto.UserRequest
		expectedData dto.VitalData
	}{
		{
			name: "1",
			userRequest: dto.UserRequest{
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
			expectedData: dto.VitalData{
				Undelivery:        5.0,  // Capacity - Tons = 10 - 5 = 5
				OperatingDistance: 400,  // Distance * 2 * QuantityTrips = 100 * 2 * 2 = 400
				Wastage:           32.0, // OperatingDistance * Consumption / 100 = 400 * 8 / 100 = 32.0
				Lifting:           3.0,  // QuantityTrips * Lifting = 2 * 1.5 = 3.0
				Underfuel:         10.0, // Undelivery * QuantityTrips * Distance / 100 = 5 * 2 * 100 / 100 = 10.0
				TotalFuel:         25.0, // Wastage + Lifting - Underfuel = 32.0 + 3.0 - 10.0 = 25.0
				DailyRate:         35.0, // FuelResidue + Refuel - TotalFuel = 20.0 + 40 - 25.0 = 35.0
				DailyRun:          900,  // SpeedometerResidue + OperatingDistance = 500 + 400 = 900
			},
		},
		{
			name: "2",
			userRequest: dto.UserRequest{
				Consumption:        10.0,
				Capacity:           12,
				FuelResidue:        25.0,
				SpeedometerResidue: 400,
				Refuel:             30,
				Distance:           80,
				QuantityTrips:      3,
				Tons:               6,
				Backload:           4,
				Lifting:            2.0,
			},
			expectedData: dto.VitalData{
				Undelivery:        0,    // Max(0, Tons + Backload - Capacity) = Max(0, 6 + 4 - 12) = 4.0
				OperatingDistance: 480,  // Distance * 2 * QuantityTrips = 80 * 2 * 3 = 480
				Wastage:           48.0, // OperatingDistance * Consumption / 100 = 480 * 10 / 100 = 48.0
				Lifting:           12.0, // 2 * QuantityTrips * Lifting = 2 * 3 * 2.0 = 12.0
				Underfuel:         0,    // Undelivery * QuantityTrips * Distance / 100 = 4 * 3 * 80 / 100 = 9.6
				TotalFuel:         60,   // Wastage + Lifting + Underfuel = 48.0 + 12.0 + 9.6 = 69.6
				DailyRate:         -5,   // FuelResidue + Refuel - TotalFuel = 25.0 + 30 - 69.6 = -14.6
				DailyRun:          880,  // SpeedometerResidue + OperatingDistance = 400 + 480 = 880
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor()
			result := processor.Calculate(tt.userRequest)

			// Проверяем каждое поле результата
			assert.Equal(t, tt.expectedData.Undelivery, result.Undelivery)
			assert.Equal(t, tt.expectedData.OperatingDistance, result.OperatingDistance)
			assert.Equal(t, tt.expectedData.Wastage, result.Wastage)
			assert.Equal(t, tt.expectedData.Lifting, result.Lifting)
			assert.Equal(t, tt.expectedData.Underfuel, result.Underfuel)
			assert.Equal(t, tt.expectedData.TotalFuel, result.TotalFuel)
			assert.Equal(t, tt.expectedData.DailyRate, result.DailyRate)
			assert.Equal(t, tt.expectedData.DailyRun, result.DailyRun)
		})
	}
}
