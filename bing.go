package fy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type bingTranslator struct{}

var bing translator = new(bingTranslator)

func (b *bingTranslator) desc() (string, string) {
	return "bing", "https://cn.bing.com/translator"
}

func BingTranslate(ctx context.Context, req Request) *Response {
	return bing.translate(ctx, req)
}

func (b *bingTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(b)

	r, data, err := sendRequest(ctx, http.MethodGet, "https://www.bing.com/translator", nil, nil)
	if err != nil {
		resp.Err = fmt.Errorf("readResp error: %v", err)
		return
	}
	timestamp, token, err := b.getTimestampAndToken(string(data))
	if err != nil {
		resp.Err = err
		return
	}

	igValue, err := b.getIGValue(string(data))
	if err != nil {
		resp.Err = err
		return
	}

	req.ToLang = b.convertLanguage(req.ToLang)
	param := url.Values{
		"fromLang": {"auto-detect"},
		"to":       {req.ToLang},
		"text":     {req.Text},
		"key":      {timestamp},
		"token":    {token},
	}

	urlStr := fmt.Sprintf("https://www.bing.com/ttranslatev3?isVertical=1&IG=%s&IID=translator.5023.1", igValue)
	body := strings.NewReader(param.Encode())
	_, data, err = sendRequest(ctx, "POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", UserAgent)
		addCookies(req, r.Cookies())
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	resp.Result = jr.Get("0.translations.0.text").String()
	if resp.Result == "" {
		resp.Err = fmt.Errorf(string(data))
	}
	return
}

func (b *bingTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh-Hans"
	}

	return l
}

func (*bingTranslator) getTimestampAndToken(dataStr string) (string, string, error) {
	// dataStr := `var params_RichTranslateHelper = [1623165977131,"Q5b9PD_0XEdOagXBcPVtlnB6ZML4958D",3600000,true];`
	result := regexp.MustCompile(`var params_RichTranslateHelper = \[(?P<result>[\s\S]+),"(?P<result1>[\s\S]+)",3600000,`).FindStringSubmatch(dataStr)
	if len(result) != 3 {
		return "", "", fmt.Errorf("cannot get timestamp and token")
	}

	return result[1], result[2], nil
}

func (*bingTranslator) getIGValue(dataStr string) (string, error) {
	// IG:"678CAA0BB18740F3A5D3EF29A6D434B0"
	igResult := regexp.MustCompile(`IG:"(?P<token>\S+)",`).FindStringSubmatch(dataStr)
	if len(igResult) != 2 {
		return "", fmt.Errorf("cannot get ig value")
	}
	return igResult[1], nil
}
