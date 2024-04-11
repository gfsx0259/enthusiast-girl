package alert

import (
	"github.com/ilyakaznacheev/cleanenv"
	"io"
)

type (
	Hook struct {
		Data `json:"data"`
	}

	Data struct {
		Metric `json:"metric_alert"`
		Text   string `json:"description_text"`
		Title  string `json:"description_title"`
		Url    string `json:"web_url"`
	}

	Metric struct {
		Projects []string `json:"projects"`
	}
)

func NewStructure(r io.Reader) *Hook {
	hook := &Hook{}

	if err := cleanenv.ParseJSON(r, hook); err != nil {
		return nil
	}

	return hook
}
