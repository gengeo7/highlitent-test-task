package questions

type QuestionDto struct {
	Text string `json:"text" validate:"required"`
}
