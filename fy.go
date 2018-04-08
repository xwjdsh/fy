package fy

var (
	Translators []Translator
)

type Request struct {
	IsChinese bool
	Text      string
}

type Response struct {
	Name   string
	Result string
	Err    error
}

type Translator interface {
	Translate(*Request) *Response
}
