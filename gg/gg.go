package gg

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
	"github.com/xwjdsh/fy"
)

type google struct{}

func init() {
	fy.Register(new(google))
}

func (s *google) Desc() (string, string, string) {
	return "gg", "google", "https://translate.google.com/"
}

func (s *google) Translate(req *fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(s)
	_, data, err := fy.ReadResp(http.Get("https://translate.google.cn"))
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v", err)
		return
	}
	vq, err := getVq(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getVq error: %v", err)
		return
	}
	tk, err := calTK(vq, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calTK error: %v", err)
		return
	}

	var from, to string
	if req.IsChinese {
		from, to = "zh-CN", "en"
	} else {
		from, to = "en", "zh-CN"
	}
	u, _ := url.Parse("https://translate.google.cn/translate_a/single")
	param := u.Query()
	param.Set("client", "t")
	param.Set("sl", from)
	param.Set("hl", "zh-CN")
	param.Set("tl", to)
	param.Set("dt", "t")
	param.Set("ie", "UTF-8")
	param.Set("oe", "UTF-8")
	param.Set("source", "btn")
	param.Set("ssel", "3")
	param.Set("tsel", "3")
	param.Set("kc", "0")
	param.Set("tk", tk)
	param.Set("q", req.Text)
	u.RawQuery = param.Encode()

	_, data, err = fy.SendRequest("GET", u.String(), nil, func(r *http.Request) error {
		r.Header.Set("User-Agent", "Paw/3.1.5 (Macintosh; OS X/10.13.2) GCDHTTPRequest")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v")
		return
	}

	jr := gjson.Parse(string(data))
	if !jr.Get("..0.0.0").Exists() {
		resp.Err = fmt.Errorf("cannot get google translate result")
		return
	}
	resp.Result = jr.Get("..0.0.0.0").String()
	return
}
