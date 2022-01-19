package fy

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type youdaoTranslator struct{}

var youdao translator = new(youdaoTranslator)

func (y *youdaoTranslator) Desc() (string, string) {
	return "youdao", "https://fanyi.youdao.com/"
}

func YoudaoTranslate(ctx context.Context, req Request) *Response {
	return youdao.translate(ctx, req)
}

func (y *youdaoTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(y)

	r, _, err := sendRequest(ctx, http.MethodGet, "https://youdao.com", nil, nil)
	if err != nil {
		resp.Err = fmt.Errorf("notReadResp error: %v", err)
		return
	}
	cookies := r.Cookies()

	// S = "fanyideskweb"
	// D = "ebSeFb%=XZ%T[KZ)c(sy!"
	// r = "" + ((new Date).getTime() + parseInt(10 * Math.random(), 10))
	// o = u.md5(S + n + r + D);
	salt := fmt.Sprintf("%d", int(time.Now().UnixNano()/1000000)+rand.Intn(10))
	h := md5.New()
	h.Write([]byte("fanyideskweb" + req.Text + salt + `Y2FYu%TNSbMCxc3t2u^XT`))
	sign := hex.EncodeToString(h.Sum(nil))

	req.ToLang = y.convertLanguage(req.ToLang)
	param := url.Values{
		"from":    {"AUTO"},
		"to":      {req.ToLang},
		"i":       {req.Text},
		"client":  {"fanyideskweb"},
		"salt":    {salt},
		"sign":    {sign},
		"version": {"3.0"},
		"keyfrom": {"fanyi.web"},
	}
	urlStr := "https://fanyi.youdao.com/translate_o"
	body := strings.NewReader(param.Encode())
	_, data, err := sendRequest(ctx, "POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Referer", "https://fanyi.youdao.com/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		addCookies(req, cookies)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("errorCode").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
		return
	}

	resp.Result = jr.Get("translateResult.0").String()
	return
}

func (y *youdaoTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh-CHS"
	}

	return l
}
