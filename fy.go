package fy

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	Chinese  = "zh-CN"
	English  = "en"
	Russian  = "ru"
	Japanese = "ja"
	German   = "de"
	French   = "fr"
	Korean   = "ko"
	Spanish  = "es"
)

var (
	translatorMap = map[string]Translator{}
	lock          sync.Mutex
)

// Register register a translator
func Register(t Translator) {
	lock.Lock()
	defer lock.Unlock()

	name, _, _ := t.Desc()
	if _, ok := translatorMap[name]; ok {
		log.Panicf("%s has been registered", name)
	}
	translatorMap[name] = t
}

// Filter filter translators
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

// NewResp new Response, set fullname
func NewResp(t Translator) *Response {
	_, fullname, _ := t.Desc()
	return &Response{
		FullName: fullname,
	}
}

// Request translate request
type Request struct {
	// The target language of translation
	TargetLang string
	// Text translate text
	Text string
}

// Response translate response
type Response struct {
	// FullName translator fullname
	FullName string
	// Result translate result
	Result string
	// Err translate error
	Err error
}

// Translator translator interface
type Translator interface {
	// Desc get translator info
	Desc() (name string, fullname string, source string)
	// Translate do translation task
	Translate(Request) *Response
}
