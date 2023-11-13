package command

type Command interface {
	Run() error
}

type ApplicationParams struct {
	Application string
	Tag         string
}
