package fy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type tencentTranslator struct{}

var tencent translator = new(tencentTranslator)

func (t *tencentTranslator) Desc() (string, string) {
	return "tencent", "https://fanyi.qq.com/"
}

func TencentTranslate(ctx context.Context, req Request) *Response {
	return tencent.translate(ctx, req)
}

func (t *tencentTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(t)

	_, data, err := sendRequest(ctx, "POST", "https://fanyi.qq.com/api/reauth12f", nil, nil)
	if err != nil {
		resp.Err = err
		return
	}
	m := map[string]string{}
	if err := json.Unmarshal(data, &m); err != nil {
		resp.Err = err
		return
	}
	qtv, qtk := m["qtv"], m["qtk"]

	req.ToLang = t.convertLanguage(req.ToLang)
	param := url.Values{
		"source":      {"auto"},
		"target":      {req.ToLang},
		"sourceText":  {req.Text},
		"qtv":         {qtv},
		"qtk":         {qtk},
		"sessionUuid": {fmt.Sprintf("translate_uuid%d", time.Now().UnixNano()/int64(time.Millisecond))},
	}

	urlStr := "https://fanyi.qq.com/api/translate"
	body := strings.NewReader(param.Encode())
	_, data, err = sendRequest(ctx, "POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Origin", "https://fanyi.qq.com")
		req.Header.Set("uc", "qEcP9NeHbGbSxJMTi6uKB7KximMaJOZIM0yohxonQs8=")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if !jr.Get("translate").Exists() {
		resp.Err = fmt.Errorf("json result translate not exists")
		return
	}
	if errorCode := jr.Get("translate.errCode").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
		return
	}

	jsonResult := jr.Get("translate.records").Array()
	for _, r := range jsonResult {
		resp.Result += r.Get("..0.targetText").String()
	}
	return
}

func (*tencentTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh"
	case Japanese:
		l = "jp"
	case Korean:
		l = "kr"
	}

	return l
}
