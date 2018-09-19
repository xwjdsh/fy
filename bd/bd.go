package bd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/xwjdsh/fy"

	"github.com/tidwall/gjson"
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
	return "bd", "baidu", "https://fanyi.baidu.com/"
}

func (b *baidu) Translate(req fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(b)

	r, err := fy.NotReadResp(http.Get("https://www.baidu.com"))
	if err != nil {
		resp.Err = fmt.Errorf("fy.NotReadResp error: %v", err)
		return
	}
	baiduCookies := r.Cookies()

	r, err = fy.NotReadResp(http.Get("https://fanyi.baidu.com"))
	if err != nil {
		resp.Err = fmt.Errorf("fy.NotReadResp error: %v", err)
		return
	}
	baiduFanyiCookies := r.Cookies()

	param := url.Values{"query": {req.Text}}
	detectUrl := "https://fanyi.baidu.com/langdetect"
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

	r, data, err = fy.SendRequest("GET", "https://fanyi.baidu.com", nil, func(req *http.Request) error {
		fy.AddCookies(req, baiduCookies)
		fy.AddCookies(req, baiduFanyiCookies)
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
	urlStr := "https://fanyi.baidu.com/v2transapi"
	body = strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		fy.AddCookies(req, baiduCookies)
		fy.AddCookies(req, baiduFanyiCookies)
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
