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

func init() {
	fy.Register(new(baidu))
}

func (b *baidu) Desc() (string, string, string) {
	return "bd", "baidu", "http://fanyi.baidu.com/"
}

func (b *baidu) Translate(req *fy.Request) (resp *fy.Response) {
	resp = &fy.Response{}
	var from, to string
	if req.IsChinese {
		from, to = "zh", "en"
	} else {
		from, to = "en", "zh"
	}

	r, err := fy.NotReadResp(http.Get("http://www.baidu.com"))
	if err != nil {
		resp.Err = fmt.Errorf("fy.NotReadResp error: %v\n", err)
		return
	}
	cookies := r.Cookies()

	r, data, err := fy.SendRequest("GET", "http://fanyi.baidu.com", nil, func(req *http.Request) error {
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		err = fmt.Errorf("fy.SendRequest error: %v\n", err)
		return
	}

	token, gtk, err := getTokenAndGtk(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getTokenAndGtk error: %v\n", err)
		return
	}
	sign, err := getSign(gtk, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("getSign error: %v\n", err)
		return
	}
	param := url.Values{
		"from":  {from},
		"to":    {to},
		"query": {req.Text},
		"sign":  {sign},
		"token": {token},
	}
	urlStr := "http://fanyi.baidu.com/v2transapi"
	body := strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.SendRequest error: %v\n", err)
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("error").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result error is %d", errorCode)
		return
	}

	_, fullname, _ := b.Desc()
	resp.FullName = fullname
	resp.Result = jr.Get("trans_result.data.0.dst").String()
	return
}
