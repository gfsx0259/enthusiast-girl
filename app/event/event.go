package event

type Event struct {
	ChatId       int64
	FromId       int64
	FromUsername string
	Message      string
}
