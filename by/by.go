package qq

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/xwjdsh/fy"
)

type bing struct{}

func init() {
	fy.Register(new(bing))
}

func (b *bing) Desc() (string, string, string) {
	return "by", "bing", "https://cn.bing.com/translator/"
}

func (b *bing) Translate(req *fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(b)

	var from, to string
	if req.IsChinese {
		from, to = "zh-CHS", "en"
	} else {
		from, to = "en", "zh-CHS"
	}
	_, data, err := fy.SendRequest("GET", "https://cn.bing.com/translator/", nil, nil)
	if err != nil {
		err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	param := url.Values{
		"from": {from},
		"to":   {to},
		"text": {req.Text},
	}

	urlStr := "https://cn.bing.com/ttranslate?IG="
	body := strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("statusCode").Int(); errorCode != 200 {
		resp.Err = fmt.Errorf("json result statusCodeis %d", errorCode)
		return
	}

	resp.Result = jr.Get("translationResponse").String()
	return
}
