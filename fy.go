package fy

import (
	"context"
	"sync"
	"time"
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

	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"
)

var translators = []translator{
	baidu, bing, google, sogou, tencent, youdao,
}

// Request translation request
type Request struct {
	// FromLang the from language of the translation
	FromLang string
	// ToLang the to language of the translation
	ToLang string
	// Text translation text
	Text string
}

// Response translation response
type Response struct {
	// Name the translator name
	Name string
	// Homepage the translator homepage
	Homepage string
	// Result the translation result
	Result string
	// Err the translation error
	Err error
}

func newResp(t translator) *Response {
	name, homepage := t.desc()
	return &Response{
		Name:     name,
		Homepage: homepage,
	}
}

type translator interface {
	desc() (name string, source string)
	translate(context.Context, Request) *Response
}

func AsyncTranslate(eachTranslatorTimeout time.Duration, req *Request, ts ...string) <-chan *Response {
	var limitMap map[string]bool
	if len(ts) > 0 {
		limitMap = map[string]bool{}
		for _, name := range ts {
			limitMap[name] = true
		}
	}

	wg := sync.WaitGroup{}
	ch := make(chan *Response)
	for _, t := range translators {
		name, _ := t.desc()
		if limitMap != nil && !limitMap[name] {
			continue
		}

		wg.Add(1)
		go func(t translator) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), eachTranslatorTimeout)
			defer cancel()
			ch <- t.translate(ctx, *req)

		}(t)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
