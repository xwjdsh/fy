package by

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/xwjdsh/fy"
)

type bing struct{}

var langConvertMap = map[string]string{
	fy.Chinese: "zh-CHS",
}

func init() {
	fy.Register(new(bing))
}

func (b *bing) Desc() (string, string, string) {
	return "by", "bing", "https://www.bing.com/translator/"
}

func (b *bing) Translate(req fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(b)

	detectUrl := "https://www.bing.com/tdetect/"
	param := url.Values{"text": {req.Text}}
	body := strings.NewReader(param.Encode())
	_, data, err := fy.SendRequest("POST", detectUrl, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	from := string(data)

	if tl, ok := langConvertMap[req.TargetLang]; ok {
		req.TargetLang = tl
	}
	param = url.Values{
		"from": {from},
		"to":   {req.TargetLang},
		"text": {req.Text},
	}

	urlStr := "https://www.bing.com/ttranslate/"
	body = strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("statusCode").Int(); errorCode != 200 {
		resp.Err = fmt.Errorf("json result statusCodeis %d", errorCode)
		return
	}

	resp.Result = jr.Get("translationResponse").String()
	return
}
