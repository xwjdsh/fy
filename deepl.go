package fy

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type deeplTranslator struct{}

var deepl translator = new(deeplTranslator)

func (b *deeplTranslator) Desc() (string, string) {
	return "deepl", "https://www.deepl.com/translator"
}

func DeeplTranslate(ctx context.Context, req Request) *Response {
	return deepl.translate(ctx, req)
}

func (t *deeplTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(t)
	req.ToLang, resp.Err = t.convertLanguage(req.ToLang)
	if resp.Err != nil {
		return
	}
	jsonStr := fmt.Sprintf(`{"jsonrpc":"2.0","method": "LMT_handle_jobs","params":{"jobs":[{"kind":"default","sentences":[{"text":"%s","id":0,"prefix":""}],"raw_en_context_before":[],"raw_en_context_after":[],"preferred_num_beams":4}],"lang":{"preference":{"weight":{},"default":"default"},"source_lang_computed":"%s","target_lang":"%s"},"priority":1,"commonJobParams":{"browserType":1},"timestamp": %d}}`, req.Text, "auto", req.ToLang, time.Now().UnixMilli())

	_, data, err := sendRequest(ctx, "POST", "https://www2.deepl.com/jsonrpc?method=LMT_handle_jobs", strings.NewReader(jsonStr), func(r *http.Request) error {
		r.Header.Set("user-agent", UserAgent)
		r.Header.Set("content-type", "application/json")
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	dit := jr.Get("result.translations.0.beams.0.sentences.0.text").String()
	if dit == "" {
		resp.Err = fmt.Errorf("cannot get translate result, data: %s", string(data))
	} else {
		resp.Result = dit
	}

	return
}

func (*deeplTranslator) convertLanguage(language string) (string, error) {
	var l string
	switch language {
	case Chinese:
		l = "ZH"
	case Japanese:
		l = "JA"
	case English:
		l = "EN"
	case Russian:
		l = "RU"
	case German:
		l = "DE"
	case French:
		l = "FR"
	case Spanish:
		l = "ES"
	}
	if l == "" {
		return "", fmt.Errorf("unsupported target language: %s", language)
	}

	return l, nil
}
