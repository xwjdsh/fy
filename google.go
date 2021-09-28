package fy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type googleTranslator struct{}

var google translator = new(googleTranslator)

func (*googleTranslator) desc() (string, string) {
	return "google", "https://translate.google.cn/"
}

func GoogleTranslate(ctx context.Context, req Request) *Response {
	return google.translate(ctx, req)
}

func (g *googleTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(g)
	reqStr := fmt.Sprintf(`[[["MkEWBc","[[\"%s\",\"%s\",\"%s\",true],[null]]",null,"generic"]]]`, req.Text, "auto", req.ToLang)
	param := url.Values{
		"f.req": {reqStr},
	}
	body := strings.NewReader(param.Encode())

	_, data, err := sendRequest(ctx, "POST", "https://translate.google.cn/_/TranslateWebserverUi/data/batchexecute", body, func(r *http.Request) error {
		r.Header.Set("user-agent", UserAgent)
		r.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	resp.Result, resp.Err = g.getResult(string(data))
	return
}

func (g *googleTranslator) getResult(dataStr string) (string, error) {
	result := regexp.MustCompile(`(?U)null,\[\[\\"(?P<result>.+)\\",null,null,null`).FindStringSubmatch(dataStr)
	if len(result) != 2 {
		return "", fmt.Errorf("cannot get result")
	}
	return result[1], nil
}
