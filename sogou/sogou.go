package sogou

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/xwjdsh/fy"
)

type Sogou struct{}

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
	fy.Translators = append(fy.Translators, &Sogou{})
}

func (sg Sogou) Name() string {
	return "sogou"
}

func (sg *Sogou) Translate(isChinese bool, text string) (string, error) {
	var from, to string
	if isChinese {
		from, to = "zh-CHS", "en"
	} else {
		from, to = "en", "zh-CHS"
	}
	param := url.Values{
		"from": {from},
		"to":   {to},
		"text": {text},
	}
	resp, err := http.PostForm("https://fanyi.sogou.com/reventondc/translate", param)
	if err != nil {
		return "", errors.Wrap(err, "http.PostForm error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "ioutil.ReadAll error")
	}
	r := &result{}
	if err := json.Unmarshal(body, r); err != nil {
		return "", errors.Wrap(err, "json.Unmarshal error")
	}
	if r.ErrorCode != 0 {
		return "", fmt.Errorf("r.ErrorCode is %d", r.ErrorCode)
	}
	if r.Translate.ErrorCode != "0" {
		return "", fmt.Errorf("r.Translate.ErrorCode is %s", r.Translate.ErrorCode)
	}

	return r.Translate.Dit, nil
}
