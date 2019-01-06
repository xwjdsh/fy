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
	return "sg", "sogou", "https://fanyi.sogou.com"
}

func (s *sogou) Translate(req fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(s)

	if tl, ok := langConvertMap[req.TargetLang]; ok {
		req.TargetLang = tl
	}

	sign, err := calSign("auto", req.TargetLang, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calSign error: %v", err)
		return
	}

	param := url.Values{
		"from": {"auto"},
		"to":   {req.TargetLang},
		"text": {req.Text},
		"s":    {sign},
	}
	urlStr := "https://fanyi.sogou.com/reventondc/translateV1"
	_, data, err := fy.ReadResp(http.PostForm(urlStr, param))
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))

	dit := jr.Get("data.translate.dit").String()
	if dit == "" {
		resp.Err = fmt.Errorf("cannot get translate result")
	} else {
		resp.Result = dit
	}

	return
}
