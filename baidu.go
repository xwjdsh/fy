package fy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
)

type baiduTranslator struct{}

var baidu translator = new(baiduTranslator)

func (*baiduTranslator) desc() (string, string) {
	return "baidu", "https://fanyi.baidu.com/"
}

func BaiduTranslate(ctx context.Context, req Request) (resp *Response) {
	return baidu.translate(ctx, req)
}

func (b *baiduTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(b)

	r, _, err := sendRequest(ctx, http.MethodGet, "https://www.baidu.com", nil, nil)
	if err != nil {
		resp.Err = fmt.Errorf("notReadResp error: %v", err)
		return
	}
	cookies := r.Cookies()

	r, _, err = sendRequest(ctx, http.MethodGet, "https://fanyi.baidu.com", nil, nil)
	if err != nil {
		resp.Err = fmt.Errorf("notReadResp error: %v", err)
		return
	}
	fanyiCookies := r.Cookies()

	param := url.Values{"query": {req.Text}}
	detectUrl := "https://fanyi.baidu.com/langdetect"
	body := strings.NewReader(param.Encode())
	_, data, err := sendRequest(ctx, http.MethodPost, detectUrl, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		return nil
	})
	if err != nil {
		resp.Err = err
		return
	}

	jr := gjson.Parse(string(data))
	if errorCode := jr.Get("error").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("langdetect json result error is %d", errorCode)
		return
	}
	from := jr.Get("lan").String()

	_, data, err = sendRequest(ctx, http.MethodGet, "https://fanyi.baidu.com", nil, func(req *http.Request) error {
		addCookies(req, cookies)
		addCookies(req, fanyiCookies)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("SendRequest error: %v", err)
		return
	}

	token, gtk, err := b.tokenAndGtk(string(data))
	if err != nil {
		resp.Err = fmt.Errorf("getTokenAndGtk error: %v", err)
		return
	}
	sign, err := b.sign(gtk, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calSign error: %v", err)
		return
	}
	req.ToLang = b.convertLanguage(req.ToLang)
	param = url.Values{
		"from":  {from},
		"to":    {req.ToLang},
		"query": {req.Text},
		"sign":  {sign},
		"token": {token},
	}
	urlStr := "https://fanyi.baidu.com/v2transapi"
	body = strings.NewReader(param.Encode())
	_, data, err = sendRequest(ctx, "POST", urlStr, body, func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		addCookies(req, cookies)
		addCookies(req, fanyiCookies)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("sendRequest error: %v", err)
		return
	}

	jr = gjson.Parse(string(data))
	if errorCode := jr.Get("error").Int(); errorCode != 0 {
		resp.Err = fmt.Errorf("json result error is %d", errorCode)
		return
	}

	resp.Result = jr.Get("trans_result.data.0.dst").String()
	return
}

func (*baiduTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh"
	case Korean:
		l = "kor"
	case Japanese:
		l = "jp"
	case French:
		l = "fra"
	case Spanish:
		l = "spa"
	}

	return l
}

const baiduSignJS = `
    function a(r) {
        if (Array.isArray(r)) {
            for (var o = 0, t = Array(r.length); o < r.length; o++) t[o] = r[o];
            return t
        }
        return Array.from(r)
    }

    function n(r, o) {
        for (var t = 0; t < o.length - 2; t += 3) {
            var a = o.charAt(t + 2);
            a = a >= "a" ? a.charCodeAt(0) - 87 : Number(a), a = "+" === o.charAt(t + 1) ? r >>> a : r << a, r = "+" === o.charAt(t) ? r + a & 4294967295 : r ^ a
        }
        return r
    }

    function e(r) {
				var i = null;
        var o = r.match(/[\uD800-\uDBFF][\uDC00-\uDFFF]/g);
        if (null === o) {
            var t = r.length;
            t > 30 && (r = "" + r.substr(0, 10) + r.substr(Math.floor(t / 2) - 5, 10) + r.substr(-10, 10))
        } else {
            for (var e = r.split(/[\uD800-\uDBFF][\uDC00-\uDFFF]/), C = 0, h = e.length, f = []; h > C; C++) "" !== e[C] && f.push.apply(f, a(e[C].split(""))), C !== h - 1 && f.push(o[C]);
            var g = f.length;
            g > 30 && (r = f.slice(0, 10).join("") + f.slice(Math.floor(g / 2) - 5, Math.floor(g / 2) + 5).join("") + f.slice(-10).join(""))
        }
        var u = void 0,
            l = "" + String.fromCharCode(103) + String.fromCharCode(116) + String.fromCharCode(107);
        u = null !== i ? i : (i = gtk || "") || "";
        for (var d = u.split("."), m = Number(d[0]) || 0, s = Number(d[1]) || 0, S = [], c = 0, v = 0; v < r.length; v++) {
            var A = r.charCodeAt(v);
            128 > A ? S[c++] = A : (2048 > A ? S[c++] = A >> 6 | 192 : (55296 === (64512 & A) && v + 1 < r.length && 56320 === (64512 & r.charCodeAt(v + 1)) ? (A = 65536 + ((1023 & A) << 10) + (1023 & r.charCodeAt(++v)), S[c++] = A >> 18 | 240, S[c++] = A >> 12 & 63 | 128) : S[c++] = A >> 12 | 224, S[c++] = A >> 6 & 63 | 128), S[c++] = 63 & A | 128)
        }
        for (var p = m, F = "" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(97) + ("" + String.fromCharCode(94) + String.fromCharCode(43) + String.fromCharCode(54)), D = "" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(51) + ("" + String.fromCharCode(94) + String.fromCharCode(43) + String.fromCharCode(98)) + ("" + String.fromCharCode(43) + String.fromCharCode(45) + String.fromCharCode(102)), b = 0; b < S.length; b++) p += S[b], p = n(p, F);
        return p = n(p, D), p ^= s, 0 > p && (p = (2147483647 & p) + 2147483648), p %= 1e6, p.toString() + "." + (p ^ m)
    }
		result = e(query)
`

func (*baiduTranslator) sign(gtk, query string) (string, error) {
	vm := otto.New()
	if err := vm.Set("gtk", gtk); err != nil {
		return "", fmt.Errorf("vm.Set gtk error: %v", err)
	}
	if err := vm.Set("query", query); err != nil {
		return "", fmt.Errorf("vm.Set query error: %v", err)
	}
	value, err := vm.Run(baiduSignJS)
	if err != nil {
		return "", fmt.Errorf("vm.Run error: %v", err)
	}
	result, err := value.ToString()
	if err != nil {
		return "", fmt.Errorf("vlue.ToString error: %v", err)
	}
	return result, nil
}

func (*baiduTranslator) tokenAndGtk(dataStr string) (token, gtk string, err error) {
	tokenResult := regexp.MustCompile(`token: '(?P<token>\S+)',`).FindStringSubmatch(dataStr)
	if len(tokenResult) != 2 {
		err = fmt.Errorf("cannot get token")
		return
	}
	token = tokenResult[1]

	gtkResult := regexp.MustCompile(`window.gtk = '(?P<gtk>\S+)';`).FindStringSubmatch(dataStr)
	if len(gtkResult) != 2 {
		err = fmt.Errorf("cannot get gtk")
		return
	}
	gtk = gtkResult[1]

	return
}
