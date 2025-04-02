package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetGeneral() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏️Посчитать📝"),
			tgbotapi.NewKeyboardButton("🫡Выдать данные📂"),
			tgbotapi.NewKeyboardButton("🗑Удалить мои данные!"),
		),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}

func GetCancel() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("😬Отмена⚠️"),
		),
	)

	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}
