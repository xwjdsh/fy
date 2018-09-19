package qq

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/xwjdsh/fy"
)

type tencent struct{}

var langConvertMap = map[string]string{
	fy.Chinese:  "zh",
	fy.Japanese: "jp",
	fy.Korean:   "kr",
}

func init() {
	fy.Register(new(tencent))
}

func (t *tencent) Desc() (string, string, string) {
	return "qq", "tencent", "https://fanyi.qq.com/"
}

func (t *tencent) Translate(req fy.Request) (resp *fy.Response) {
	resp = fy.NewResp(t)

	_, data, err := fy.SendRequest("GET", "https://fanyi.qq.com", nil, nil)
	if err != nil {
		err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	qtv, qtk, err := getQtk(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getQtk error: %v", err)
		return
	}

	if tl, ok := langConvertMap[req.TargetLang]; ok {
		req.TargetLang = tl
	}
	param := url.Values{
		"source":     {"auto"},
		"target":     {req.TargetLang},
		"sourceText": {req.Text},
	}

	urlStr := "https://fanyi.qq.com/api/translate"
	body := strings.NewReader(param.Encode())
	_, data, err = fy.SendRequest("POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Origin", "http://fanyi.qq.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Cookie", fmt.Sprintf("qtv=%s; qtk=%s", qtv, qtk))
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("fy.SendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	if !jr.Get("translate").Exists() {
		resp.Err = fmt.Errorf("json result translate not exists")
		return
	}
	if errorCode := jr.Get("translate.errCode").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result errorCode is %d", errorCode)
		return
	}

	jsonResult := jr.Get("translate.records").Array()
	for _, r := range jsonResult {
		resp.Result += r.Get("..0.targetText").String()
	}
	return
}

func getQtk(dataStr string) (qtv string, qtk string, err error) {
	//document.cookie = "qtv=ad15088b8bcde724";
	qtvResult := regexp.MustCompile(`"qtv=(?P<qtv>\S+)";`).FindStringSubmatch(dataStr)
	if len(qtvResult) != 2 {
		err = fmt.Errorf("cannot get qtv")
		return
	}
	qtv = qtvResult[1]

	//document.cookie = "qtk=aK4qrfL4bLogktVEfIMc785lhWKxOuLuOF243HgKs/lOcPqPhoiwsR+7ysGoTF/rqx1EABKUpNJq2OqbE1PY9T9ICiU2Qm2l0yIMqg3mworcjCX8tiaZzYjkQQqSTk7gCIz/WY0NhTJUrrOemb6nRQ==";
	qtkResult := regexp.MustCompile(`"qtk=(?P<qtk>\S+)";`).FindStringSubmatch(dataStr)
	if len(qtkResult) != 2 {
		err = fmt.Errorf("cannot get qtk")
		return
	}
	qtk = qtkResult[1]
	return
}
