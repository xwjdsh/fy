package fy

import (
	"log"
	"sync"
)

var (
	TranslatorMap = map[string]Translator{}
	lock          sync.Mutex
)

func Register(t Translator) {
	lock.Lock()
	defer lock.Unlock()

	name, _, _ := t.Desc()
	if _, ok := TranslatorMap[name]; ok {
		log.Printf("%s has been registered", name)
	}
	TranslatorMap[name] = t
}

func NewResp(t Translator) *Response {
	_, fullname, _ := t.Desc()
	return &Response{
		FullName: fullname,
	}
}

type Request struct {
	IsChinese bool
	Text      string
}

type Response struct {
	FullName string
	Result   string
	Err      error
}

type Translator interface {
	Desc() (name string, fullname string, source string)
	Translate(*Request) *Response
}
