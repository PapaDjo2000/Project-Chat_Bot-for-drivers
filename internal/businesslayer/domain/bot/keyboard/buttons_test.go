package keyboard

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
func Test_GetGeneral(t *testing.T) {
	keyboard := GetGeneral()
	assert.IsType(t, tgbotapi.ReplyKeyboardMarkup{}, keyboard)
	assert.Len(t, keyboard.Keyboard, 1, "Keyboard should have exactly 1 row")
	assert.Len(t, keyboard.Keyboard[0], 3, "First row should have 3 buttons")
	assert.Equal(t, "✏️Посчитать📝", keyboard.Keyboard[0][0].Text)
	assert.Equal(t, "🫡Выдать данные📂", keyboard.Keyboard[0][1].Text)
	assert.Equal(t, "🗑Удалить мои данные!", keyboard.Keyboard[0][2].Text)

	assert.True(t, keyboard.ResizeKeyboard, "ResizeKeyboard should be true")
	assert.False(t, keyboard.OneTimeKeyboard, "OneTimeKeyboard should be false")
}

func Test_GetCancel(t *testing.T) {
	keyboard := GetCancel()
	assert.IsType(t, tgbotapi.ReplyKeyboardMarkup{}, keyboard)
	assert.Len(t, keyboard.Keyboard, 1, "Keyboard should have exactly 1 row")
	assert.Len(t, keyboard.Keyboard[0], 1, "First row should have 1 button")
	assert.Equal(t, "😬Отмена⚠️", keyboard.Keyboard[0][0].Text)
	assert.True(t, keyboard.ResizeKeyboard, "ResizeKeyboard should be true")
	assert.False(t, keyboard.OneTimeKeyboard, "OneTimeKeyboard should be false")
}
