package command

import "strings"

type Command interface {
	Run() error
}

type ApplicationParams struct {
	Application string
	Tag         string
}

func ResolveFinalTag(tag string) string {
	return tag[:strings.IndexByte(tag, '-')]
}
