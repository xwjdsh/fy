package sg

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
	"github.com/xwjdsh/fy"
)

type sogou struct{}

func init() {
	fy.Register(new(sogou))
}

var langConvertMap = map[string]string{
	fy.Chinese: "zh-CHS",
}

func (s *sogou) Desc() (string, string, string) {
	return "sg", "sogou", "http://fanyi.sogou.com/"
}

func (s *sogou) Translate(req *fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(s)

	if tl, ok := langConvertMap[req.TargetLang]; ok {
		req.TargetLang = tl
	}
	param := url.Values{
		"from": {"auto"},
		"to":   {req.TargetLang},
		"text": {req.Text},
	}
	urlStr := "https://fanyi.sogou.com/reventondc/translate"
	_, data, err := fy.ReadResp(http.PostForm(urlStr, param))
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("errorCode").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
		return
	}

	if errorCode := jr.Get("translate.errorCode").String(); errorCode != "0" {
		resp.Err = fmt.Errorf("json result translate.errorCode is %s", errorCode)
		return
	}
	resp.Result = jr.Get("translate.dit").String()
	return
}
