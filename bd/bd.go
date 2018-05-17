package bd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/xwjdsh/fy"
)

type baidu struct{}

var langConvertMap = map[string]string{
	fy.Chinese:  "zh",
	fy.Korean:   "kor",
	fy.Japanese: "jp",
	fy.French:   "fra",
	fy.Spanish:  "spa",
}

func init() {
	fy.Register(new(baidu))
}

func (b *baidu) Desc() (string, string, string) {
	return "bd", "baidu", "http://fanyi.baidu.com/"
}

func (b *baidu) Translate(req fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(b)

	r, err := fy.NotReadResp(http.Get("http://www.baidu.com"))
	if err != nil {
		resp.Err = fmt.Errorf("fy.NotReadResp error: %v", err)
		return
	}
	cookies := r.Cookies()

	param := url.Values{"query": {req.Text}}
	detectUrl := "http://fanyi.baidu.com/langdetect"
	body := strings.NewReader(param.Encode())
	r, data, err := fy.SendRequest("POST", detectUrl, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		return nil
	})

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("error").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("langdetect json result error is %d", errorCode)
		return
	}
	from := jr.Get("lan").String()

	r, data, err = fy.SendRequest("GET", "http://fanyi.baidu.com", nil, func(req *http.Request) error {
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	token, gtk, err := getTokenAndGtk(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getTokenAndGtk error: %v", err)
		return
	}
	sign, err := calSign(gtk, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calSign error: %v", err)
		return
	}
	if tl, ok := langConvertMap[req.TargetLang]; ok {
		req.TargetLang = tl
	}
	param = url.Values{
		"from":  {from},
		"to":    {req.TargetLang},
		"query": {req.Text},
		"sign":  {sign},
		"token": {token},
	}
	urlStr := "http://fanyi.baidu.com/v2transapi"
	body = strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	jr = gjson.Parse(string(data))
	if errorCode := jr.Get("error").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result error is %d", errorCode)
		return
	}

	resp.Result = jr.Get("trans_result.data.0.dst").String()
	return
}
