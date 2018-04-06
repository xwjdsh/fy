package bd

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/xwjdsh/fy"
)

type baidu struct{}

func init() {
	fy.Translators = append(fy.Translators, &baidu{})
}

func (b *baidu) Translate(req *fy.Request) (resp *fy.Response) {
	resp = &fy.Response{}
	var from, to string
	if req.IsChinese {
		from, to = "zh", "en"
	} else {
		from, to = "en", "zh"
	}

	r, err := fy.NotReadResp(http.Get("http://fanyi.baidu.com"))
	if err != nil {
		resp.Err = errors.Wrap(err, "fy.NotReadResp error")
		return
	}
	cookies := r.Cookies()
	for _, cookie := range cookies {
		log.Println(cookie.Name, cookie.Value)

	}

	r, data, err := fy.SendRequest("GET", "http://fanyi.baidu.com", nil, func(req *http.Request) error {
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "fy.SendRequest error")
		return
	}

	token, gtk, err := getTokenAndGtk(string(data))
	if err != nil {
		resp.Err = errors.Wrap(err, "getTokenAndGtk error")
		return
	}
	sign, err := getSign(gtk, req.Text)
	if err != nil {
		resp.Err = errors.Wrap(err, "getSign error")
		return
	}
	param := url.Values{
		"from":              {from},
		"to":                {to},
		"query":             {req.Text},
		"transtype":         {"realtime"},
		"simple_means_flag": {"3"},
		"sign":              {sign},
		"token":             {token},
	}
	urlStr := "http://fanyi.baidu.com/v2transapi"
	body := strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		fy.AddCookies(req, cookies)
		return nil
	})
	if err != nil {
		resp.Err = errors.Wrap(err, "fy.SendRequest error")
		return
	}
	resp.Result = string(data)
	return
}
