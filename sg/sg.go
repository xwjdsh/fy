package sg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/xwjdsh/fy"
)

type sogou struct{}

type result struct {
	ErrorCode int `json:"errorCode"`
	Detect    struct {
		Zly       string `json:"zly"`
		IsCache   bool   `json:"is_cache"`
		Detect    string `json:"detect"`
		ErrorCode string `json:"errorCode"`
		Language  string `json:"language"`
		ID        string `json:"id"`
		Text      string `json:"text"`
	} `json:"detect"`
	Message   string `json:"message"`
	Translate struct {
		QcType    string `json:"qc_type"`
		Zly       string `json:"zly"`
		IsCache   bool   `json:"is_cache"`
		ErrorCode string `json:"errorCode"`
		Index     string `json:"index"`
		Source    string `json:"source"`
		Dit       string `json:"dit"`
		From      string `json:"from"`
		Text      string `json:"text"`
		To        string `json:"to"`
		ID        string `json:"id"`
		OrigText  string `json:"orig_text"`
		Md5       string `json:"md5"`
	} `json:"translate"`
	IsHasOxford  bool `json:"isHasOxford"`
	IsHasChinese bool `json:"isHasChinese"`
}

func init() {
	fy.Translators = append(fy.Translators, &sogou{})
}

func (s *sogou) Translate(req *fy.Request) (resp *fy.Response) {
	resp = &fy.Response{}
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
	r, err := http.PostForm("https://fanyi.sogou.com/reventondc/translate", param)
	if err != nil {
		resp.Err = errors.Wrap(err, "http.PostForm error")
		return
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Err = errors.Wrap(err, "ioutil.ReadAll error")
		return
	}

	rt := &result{}
	if err := json.Unmarshal(body, rt); err != nil {
		resp.Err = errors.Wrap(err, "json.Unmarshal error")
		return
	}
	if rt.ErrorCode != 0 {
		resp.Err = fmt.Errorf("rt.ErrorCode is %d", rt.ErrorCode)
		return
	}
	if rt.Translate.ErrorCode != "0" {
		resp.Err = fmt.Errorf("rt.Translate.ErrorCode is %s", rt.Translate.ErrorCode)
		return
	}

	resp.Name = "sogou"
	resp.Result = rt.Translate.Dit
	return
}
