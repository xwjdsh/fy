package fy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type bingTranslator struct{}

var bing translator = new(bingTranslator)

func (b *bingTranslator) desc() (string, string) {
	return "bing", "https://www.bing.com/translator/"
}

func BingTranslate(ctx context.Context, req Request) *Response {
	return bing.translate(ctx, req)
}

func (b *bingTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(b)

	detectUrl := "https://www.bing.com/tdetect/"
	param := url.Values{"text": {req.Text}}
	body := strings.NewReader(param.Encode())
	_, data, err := sendRequest(ctx, http.MethodPost, detectUrl, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	from := string(data)

	req.ToLang = b.convertLanguage(req.ToLang)
	param = url.Values{
		"from": {from},
		"to":   {req.ToLang},
		"text": {req.Text},
	}

	urlStr := "https://www.bing.com/ttranslate/"
	body = strings.NewReader(param.Encode())
	_, data, err = sendRequest(ctx, "POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	resp.Result = jr.Get("translationResponse").String()
	return
}

func (b *bingTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh-CHS"
	}

	return l
}
