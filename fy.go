package fy

import (
	"context"
	"unicode"

	"github.com/chromedp/chromedp"
)

func New() *Client {
	tm := map[string]translator{}
	for _, t := range []translator{new(baiduTranslator), new(bingTranslator)} {
		tm[t.name()] = t
	}

	return &Client{tm: tm}
}

type Client struct {
	tm map[string]translator
}

func (c *Client) Baidu(ctx context.Context, text string) *Response {
	return c.run(ctx, c.tm[BAIDU], text, isChinese(text))
}

func (c *Client) Bing(ctx context.Context, text string) *Response {
	return c.run(ctx, c.tm[BING], text, isChinese(text))
}

func (c *Client) run(ctx context.Context, t translator, text string, isChinese bool) *Response {
	if t == nil {
		return nil
	}
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	ctx, cancel := chromedp.NewExecAllocator(ctx, options...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	actions, result := t.getActions(text, isChinese)
	err := chromedp.Run(ctx, actions...)
	return &Response{
		Name:     t.name(),
		Homepage: t.homepage(),
		Err:      err,
		Result:   *result,
	}
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

func isChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

//
//func AsyncTranslate(eachTranslatorTimeout time.Duration, req *Request, ts ...string) <-chan *Response {
//	var limitMap map[string]bool
//	if len(ts) > 0 {
//		limitMap = map[string]bool{}
//		for _, name := range ts {
//			limitMap[name] = true
//		}
//	}
//
//	wg := sync.WaitGroup{}
//	ch := make(chan *Response)
//	for _, t := range translators {
//		name, _ := t.desc()
//		if limitMap != nil && !limitMap[name] {
//			continue
//		}
//
//		wg.Add(1)
//		go func(t translator) {
//			defer wg.Done()
//
//			ctx, cancel := context.WithTimeout(context.Background(), eachTranslatorTimeout)
//			defer cancel()
//			ch <- t.translate(ctx, *req)
//
//		}(t)
//	}
//
//	go func() {
//		wg.Wait()
//		close(ch)
//	}()
//
//	return ch
//}
