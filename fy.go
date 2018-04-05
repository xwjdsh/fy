package fy

var (
	Translators []Translator
	ResultChan  chan *Result
)

type Translator interface {
	Name() string
	Translate(bool, string) (string, error)
}

type Result struct {
	Name string
	T    string
	Err  error
}
