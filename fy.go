package fy

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	translatorMap = map[string]Translator{}
	lock          sync.Mutex
)

func Register(t Translator) {
	lock.Lock()
	defer lock.Unlock()

	name, _, _ := t.Desc()
	if _, ok := translatorMap[name]; ok {
		log.Panicf("%s has been registered", name)
	}
	translatorMap[name] = t
}
func Filter(only, except string) ([]Translator, error) {
	lock.Lock()
	defer lock.Unlock()

	getMap := func(str string) (map[string]bool, error) {
		m := map[string]bool{}
		if str == "" {
			return m, nil
		}
		for _, s := range strings.Split(str, ",") {
			if _, ok := translatorMap[s]; !ok {
				return nil, fmt.Errorf("the translator [%s] does not exist", s)
			}
			m[s] = true
		}
		return m, nil
	}
	onlyMap, err := getMap(only)
	if err != nil {
		return nil, err
	}
	exceptMap, err := getMap(except)
	if err != nil {
		return nil, err
	}
	result := []Translator{}
	for k, v := range translatorMap {
		if (len(onlyMap) > 0 && !onlyMap[k]) || (len(exceptMap) > 0 && exceptMap[k]) {
			continue
		}
		result = append(result, v)
	}
	return result, nil
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
