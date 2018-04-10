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
		resp.Err = fmt.Errorf("fy.ReadResp error: %v\n", err)
		return
	}
	vq, err := getVq(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getVq error: %v\n", err)
		return
	}
	tk, err := calTK(vq, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calTK error: %v\n", err)
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
	param["dt"] = []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"}
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
		resp.Err = fmt.Errorf("fy.ReadResp error: %v\n")
		return
	}

	jr := gjson.Parse(string(data))
	//if errorCode := jr.Get("errorCode").Int(); errorCode != 0 {
	//resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
	//return
	//}

	//if errorCode := jr.Get("translate.errorCode").String(); errorCode != "0" {
	//resp.Err = fmt.Errorf("json result translate.errorCode is %s", errorCode)
	//return
	//}
	resp.Result = jr.Get("..0.0.0.0").String()
	return
}
