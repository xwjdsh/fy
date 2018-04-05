package fy

var (
	Translators  []Translator
	ResponseChan = make(chan *Response)
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

func Handle(t Translator, r *Request) {
	ResponseChan <- t.Translate(r)
}
