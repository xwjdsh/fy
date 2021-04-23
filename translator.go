package fy

import (
	"time"

	"github.com/chromedp/chromedp"
)

const (
	BAIDU = "baidu"
	BING  = "bing"
)

type translator interface {
	name() string
	homepage() string
	getActions(text string, isChinese bool) ([]chromedp.Action, *string)
}

var (
	_ translator = new(baiduTranslator)
	_ translator = new(bingTranslator)
)

type baiduTranslator struct{}

func (*baiduTranslator) name() string {
	return BAIDU
}

func (*baiduTranslator) homepage() string {
	return "https://fanyi.baidu.com/"
}

func (b *baiduTranslator) getActions(text string, isChinese bool) ([]chromedp.Action, *string) {
	result := ""
	return []chromedp.Action{
		chromedp.Navigate(b.homepage()),
		chromedp.SendKeys("#baidu_translate_input", text, chromedp.ByID),
		chromedp.Text(".target-output", &result),
	}, &result
}

type bingTranslator struct{}

func (b *bingTranslator) name() string {
	return BING
}

func (b *bingTranslator) homepage() string {
	return "https://www.bing.com/translator"
}

func (b *bingTranslator) getActions(text string, isChinese bool) ([]chromedp.Action, *string) {
	result := ""

	if isChinese {

	} else {

	}
	return []chromedp.Action{
		chromedp.Navigate(b.homepage()),
		chromedp.SendKeys("#tta_input_ta", text, chromedp.ByID),
		chromedp.Click("#tta_tgtsl", chromedp.ByID),
		chromedp.Text("#tta_output_ta", &result, chromedp.ByID),
		chromedp.Sleep(10 * time.Second),
	}, &result
}
