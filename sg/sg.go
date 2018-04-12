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

func (s *sogou) Desc() (string, string, string) {
	return "sg", "sogou", "http://fanyi.sogou.com/"
}

func (s *sogou) Translate(req *fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(s)

	var from, to string
	if req.IsChinese {
		from, to = "zh-CHS", "en"
	} else {
		from, to = "en", "zh-CHS"
	}
	param := url.Values{
		"from": {from},
		"to":   {to},
		"text": {req.Text},
	}
	urlStr := "https://fanyi.sogou.com/reventondc/translate"
	_, data, err := fy.ReadResp(http.PostForm(urlStr, param))
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v")
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
