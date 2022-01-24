package fy

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type caiyunTranslator struct {
	token      string
	seq1, seq2 string
}

var caiyun translator = &caiyunTranslator{
	token: "token:qgemv4jr1y38jyq6vhvi",
	seq1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	seq2:  "NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm",
}

func (b *caiyunTranslator) Desc() (string, string) {
	return "caiyun", "https://fanyi.caiyunapp.com"
}

func CaiyunTranslate(ctx context.Context, req Request) *Response {
	return caiyun.translate(ctx, req)
}

func (t *caiyunTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(t)
	req.ToLang, resp.Err = t.convertLanguage(req.ToLang)
	if resp.Err != nil {
		return
	}

	browserId := strings.ReplaceAll(uuid.New().String(), "-", "")
	jwtBody := strings.NewReader(fmt.Sprintf(`{"browser_id":"%s"}`, browserId))

	var jwtResponse struct {
		Rc  int    `json:"rc"`
		Jwt string `json:"jwt"`
	}

	_, data, err := sendRequest(ctx, http.MethodPost, "https://api.interpreter.caiyunai.com/v1/user/jwt/generate", jwtBody, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("X-Authorization", t.token)
		return nil
	})
	if err := json.Unmarshal(data, &jwtResponse); err != nil {
		resp.Err = fmt.Errorf("get jwt error: %v", err)
		return
	}
	if jwtResponse.Rc != 0 {
		resp.Err = fmt.Errorf("unexpected jwt response: %s", string(data))
		return
	}

	param := map[string]interface{}{
		"browser_id": browserId,
		"cached":     true,
		"detect":     true,
		"dict":       true,
		"media":      "text",
		"os_type":    "web",
		"replaced":   true,
		"request_id": "web_fanyi",
		"source":     req.Text,
		"trans_type": "auto2" + req.ToLang,
	}
	reqData, _ := json.Marshal(param)
	body := bytes.NewReader(reqData)

	_, data, err = sendRequest(ctx, http.MethodPost, "https://api.interpreter.caiyunai.com/v1/translator", body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("X-Authorization", t.token)
		req.Header.Set("T-Authorization", jwtResponse.Jwt)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	dit := jr.Get("target").String()
	if dit == "" {
		resp.Err = fmt.Errorf("cannot get translate result, resp: %s", string(data))
		return
	}

	encoded := ""
	for _, s := range dit {
		if i := strings.IndexRune(t.seq1, s); i > -1 {
			encoded += string(t.seq2[i])
		} else {
			encoded += string(s)
		}
	}

	result, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		resp.Err = fmt.Errorf("base64 decode error: %v", err)
		return
	}

	resp.Result = string(result)
	return
}

func (*caiyunTranslator) convertLanguage(language string) (string, error) {
	var l string
	switch language {
	case Chinese:
		l = "zh"
	case Japanese, English, Russian, French, Spanish:
		l = language
	}

	if l == "" {
		return "", fmt.Errorf("unsupported target language: %s", language)
	}

	return l, nil
}
