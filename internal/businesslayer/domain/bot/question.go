package bot

const (
	// TODO: rename
	defaultFinish         = -1
	finishOfFirstQuestion = -100
)

type Question struct {
	Text    string
	Options []Option
}

type Option struct {
	Text           string
	NextQuestionID int
}

var questions = map[int]Question{
	1: {
		Text: "Что будем делать?",
		Options: []Option{
			{Text: "Посчитать", NextQuestionID: 2},
			{Text: "Показать данные", NextQuestionID: 3},
		},
	},
	2: {
		Text: "Введите Расход:",
		Options: []Option{
			{Text: "", NextQuestionID: finishOfFirstQuestion},
			{NextQuestionID: 3},
		},
	},
	3: {
		Text: "Ваша данные:",
		Options: []Option{
			{Text: "", NextQuestionID: finishOfFirstQuestion},
			{NextQuestionID: 3},
		},
	},
}
