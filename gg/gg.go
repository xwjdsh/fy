package gg

import (
	"fmt"
	"log"
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
	resp = &fy.Response{}
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
	param.Set("hl", "zh-CN")
	param.Set("sl", from)
	param.Set("tl", to)
	param.Set("ie", "UTF-8")
	param.Set("oe", "UTF-8")
	param.Set("tk", tk)
	param.Set("q", req.Text)
	u.RawQuery = param.Encode()

	_, data, err = fy.ReadResp(http.Get(u.String()))
	if err != nil {
		resp.Err = fmt.Errorf("fy.ReadResp error: %v\n")
		return
	}

	log.Println(string(data))
	jr := gjson.Parse(string(data))
	//if errorCode := jr.Get("errorCode").Int(); errorCode != 0 {
	//resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
	//return
	//}

	//if errorCode := jr.Get("translate.errorCode").String(); errorCode != "0" {
	//resp.Err = fmt.Errorf("json result translate.errorCode is %s", errorCode)
	//return
	//}
	_, fullname, _ := s.Desc()
	resp.FullName = fullname
	resp.Result = jr.Get("translate.dit").String()
	return
}
