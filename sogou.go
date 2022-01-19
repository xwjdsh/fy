package fy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
)

type sogouTranslator struct{}

var sogou translator = new(sogouTranslator)

func (s *sogouTranslator) Desc() (string, string) {
	return "sogou", "https://fanyi.sogou.com"
}

func SogouTranslate(ctx context.Context, req Request) *Response {
	return sogou.translate(ctx, req)
}

func (s *sogouTranslator) translate(ctx context.Context, req Request) (resp *Response) {
	resp = newResp(s)

	r, _, err := sendRequest(ctx, http.MethodGet, "https://fanyi.sogou.com/", nil, func(r *http.Request) error {
		r.Header.Set("User-Agent", UserAgent)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("notReadResp error: %v", err)
		return
	}
	cookies := r.Cookies()

	req.ToLang = s.convertLanguage(req.ToLang)

	sign, err := s.calSign("auto", req.ToLang, req.Text)
	if err != nil {
		resp.Err = fmt.Errorf("calSign error: %v", err)
		return
	}

	param := map[string]interface{}{
		"client":   "pc",
		"exchange": false,
		"fr":       "browser_pc",
		"from":     "auto",
		"needQc":   1,
		"s":        sign,
		"text":     req.Text,
		"to":       req.ToLang,
	}
	urlStr := "https://fanyi.sogou.com/api/transpc/text/result"

	reqData, _ := json.Marshal(param)
	body := bytes.NewReader(reqData)
	_, data, err := sendRequest(ctx, http.MethodPost, urlStr, body, func(req *http.Request) error {
		addCookies(req, cookies)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UserAgent)
		return nil
	})
	if err != nil {
		resp.Err = fmt.Errorf("readResp error: %v", err)
		return
	}

	jr := gjson.Parse(string(data))
	dit := jr.Get("data.translate.dit").String()
	if dit == "" {
		resp.Err = fmt.Errorf("cannot get translate result")
	} else {
		resp.Result = dit
	}

	return
}

func (s *sogouTranslator) convertLanguage(language string) string {
	l := language
	switch language {
	case Chinese:
		l = "zh-CHS"
	}

	return l
}

const signJS = `
    var n150 = {
      "rotl": function(t, e) {
        return t << e | t >>> 32 - e
      },
      "rotr": function(t, e) {
        return t << 32 - e | t >>> e
      },
      "endian": function(t) {
        if (t.constructor == Number)
          return 16711935 & n150.rotl(t, 8) | 4278255360 & n150.rotl(t, 24);
        for (var e = 0; e < t.length; e++)
          t[e] = n150.endian(t[e]);
        return t
      },
      "randomBytes": function(t) {
        for (var e = []; t > 0; t--)
          e.push(Math.floor(256 * Math.random()));
        return e
      },
      "bytesToWords": function(t) {
        for (var e = [], n = 0, r = 0; n < t.length; n++,
          r += 8)
          e[r >>> 5] |= t[n] << 24 - r % 32;
        return e
      },
      "wordsToBytes": function(t) {
        for (var e = [], n = 0; n < 32 * t.length; n += 8)
          e.push(t[n >>> 5] >>> 24 - n % 32 & 255);
        return e
      },
      "bytesToHex": function(t) {
        for (var e = [], n = 0; n < t.length; n++)
          e.push((t[n] >>> 4).toString(16)),
          e.push((15 & t[n]).toString(16));
        return e.join("")
      },
      "hexToBytes": function(t) {
        for (var e = [], n = 0; n < t.length; n += 2)
          e.push(parseInt(t.substr(n, 2), 16));
        return e
      },
      "bytesToBase64": function(t) {
        for (var n = [], r = 0; r < t.length; r += 3)
          for (var o = t[r] << 16 | t[r + 1] << 8 | t[r + 2], i = 0; i < 4; i++)
            8 * r + 6 * i <= 8 * t.length ? n.push(e.charAt(o >>> 6 * (3 - i) & 63)) : n.push("=");
        return n.join("")
      },
      "base64ToBytes": function(t) {
        t = t.replace(/[^A-Z0-9+\/]/gi, "");
        for (var n = [], r = 0, o = 0; r < t.length; o = ++r % 4)
          0 != o && n.push((e.indexOf(t.charAt(r - 1)) & Math.pow(2, -2 * o + 8) - 1) << 2 * o | e.indexOf(t.charAt(r)) >>> 6 - 2 * o);
        return n
      }
    };

    var n128 = {
      "utf8": {
        "stringToBytes": function(t) {
          return n128.bin.stringToBytes(unescape(encodeURIComponent(t)))
        },
        "bytesToString": function(t) {
          return decodeURIComponent(escape(n128.bin.bytesToString(t)))
        }
      },
      "bin": {
        "stringToBytes": function(t) {
          for (var e = [], n = 0; n < t.length; n++)
            e.push(255 & t.charCodeAt(n));
          return e
        },
        "bytesToString": function(t) {
          for (var e = [], n = 0; n < t.length; n++)
            e.push(String.fromCharCode(t[n]));
          return e.join("")
        }
      }
    };

    function a(t, n) {
      var r = n128.utf8;
      var i = n128.bin;
      var e = n150;
      t.constructor == String ? t = n && "binary" === n.encoding ? i.stringToBytes(t) : r.stringToBytes(t) : o(t) ? t = Array.prototype.slice.call(t, 0) : Array.isArray(t) || (t = t.toString());
      for (var c = e.bytesToWords(t), u = 8 * t.length, s = 1732584193, l = -271733879, f = -1732584194, p = 271733878, h = 0; h < c.length; h++)
        c[h] = 16711935 & (c[h] << 8 | c[h] >>> 24) | 4278255360 & (c[h] << 24 | c[h] >>> 8);
      c[u >>> 5] |= 128 << u % 32,
        c[14 + (u + 64 >>> 9 << 4)] = u;
      for (var d = aff, y = agg, v = ahh, g = aii, h = 0; h < c.length; h += 16) {
        var b = s,
          m = l,
          w = f,
          T = p;
        s = d(s, l, f, p, c[h + 0], 7, -680876936),
          p = d(p, s, l, f, c[h + 1], 12, -389564586),
          f = d(f, p, s, l, c[h + 2], 17, 606105819),
          l = d(l, f, p, s, c[h + 3], 22, -1044525330),
          s = d(s, l, f, p, c[h + 4], 7, -176418897),
          p = d(p, s, l, f, c[h + 5], 12, 1200080426),
          f = d(f, p, s, l, c[h + 6], 17, -1473231341),
          l = d(l, f, p, s, c[h + 7], 22, -45705983),
          s = d(s, l, f, p, c[h + 8], 7, 1770035416),
          p = d(p, s, l, f, c[h + 9], 12, -1958414417),
          f = d(f, p, s, l, c[h + 10], 17, -42063),
          l = d(l, f, p, s, c[h + 11], 22, -1990404162),
          s = d(s, l, f, p, c[h + 12], 7, 1804603682),
          p = d(p, s, l, f, c[h + 13], 12, -40341101),
          f = d(f, p, s, l, c[h + 14], 17, -1502002290),
          l = d(l, f, p, s, c[h + 15], 22, 1236535329),
          s = y(s, l, f, p, c[h + 1], 5, -165796510),
          p = y(p, s, l, f, c[h + 6], 9, -1069501632),
          f = y(f, p, s, l, c[h + 11], 14, 643717713),
          l = y(l, f, p, s, c[h + 0], 20, -373897302),
          s = y(s, l, f, p, c[h + 5], 5, -701558691),
          p = y(p, s, l, f, c[h + 10], 9, 38016083),
          f = y(f, p, s, l, c[h + 15], 14, -660478335),
          l = y(l, f, p, s, c[h + 4], 20, -405537848),
          s = y(s, l, f, p, c[h + 9], 5, 568446438),
          p = y(p, s, l, f, c[h + 14], 9, -1019803690),
          f = y(f, p, s, l, c[h + 3], 14, -187363961),
          l = y(l, f, p, s, c[h + 8], 20, 1163531501),
          s = y(s, l, f, p, c[h + 13], 5, -1444681467),
          p = y(p, s, l, f, c[h + 2], 9, -51403784),
          f = y(f, p, s, l, c[h + 7], 14, 1735328473),
          l = y(l, f, p, s, c[h + 12], 20, -1926607734),
          s = v(s, l, f, p, c[h + 5], 4, -378558),
          p = v(p, s, l, f, c[h + 8], 11, -2022574463),
          f = v(f, p, s, l, c[h + 11], 16, 1839030562),
          l = v(l, f, p, s, c[h + 14], 23, -35309556),
          s = v(s, l, f, p, c[h + 1], 4, -1530992060),
          p = v(p, s, l, f, c[h + 4], 11, 1272893353),
          f = v(f, p, s, l, c[h + 7], 16, -155497632),
          l = v(l, f, p, s, c[h + 10], 23, -1094730640),
          s = v(s, l, f, p, c[h + 13], 4, 681279174),
          p = v(p, s, l, f, c[h + 0], 11, -358537222),
          f = v(f, p, s, l, c[h + 3], 16, -722521979),
          l = v(l, f, p, s, c[h + 6], 23, 76029189),
          s = v(s, l, f, p, c[h + 9], 4, -640364487),
          p = v(p, s, l, f, c[h + 12], 11, -421815835),
          f = v(f, p, s, l, c[h + 15], 16, 530742520),
          l = v(l, f, p, s, c[h + 2], 23, -995338651),
          s = g(s, l, f, p, c[h + 0], 6, -198630844),
          p = g(p, s, l, f, c[h + 7], 10, 1126891415),
          f = g(f, p, s, l, c[h + 14], 15, -1416354905),
          l = g(l, f, p, s, c[h + 5], 21, -57434055),
          s = g(s, l, f, p, c[h + 12], 6, 1700485571),
          p = g(p, s, l, f, c[h + 3], 10, -1894986606),
          f = g(f, p, s, l, c[h + 10], 15, -1051523),
          l = g(l, f, p, s, c[h + 1], 21, -2054922799),
          s = g(s, l, f, p, c[h + 8], 6, 1873313359),
          p = g(p, s, l, f, c[h + 15], 10, -30611744),
          f = g(f, p, s, l, c[h + 6], 15, -1560198380),
          l = g(l, f, p, s, c[h + 13], 21, 1309151649),
          s = g(s, l, f, p, c[h + 4], 6, -145523070),
          p = g(p, s, l, f, c[h + 11], 10, -1120210379),
          f = g(f, p, s, l, c[h + 2], 15, 718787259),
          l = g(l, f, p, s, c[h + 9], 21, -343485551),
          s = s + b >>> 0,
          l = l + m >>> 0,
          f = f + w >>> 0,
          p = p + T >>> 0
      }
      return e.endian([s, l, f, p])
    };

    function aff(t, e, n, r, o, i, a) {
      var c = t + (e & n | ~e & r) + (o >>> 0) + a;
      return (c << i | c >>> 32 - i) + e
    }

    function agg(t, e, n, r, o, i, a) {
      var c = t + (e & r | n & ~r) + (o >>> 0) + a;
      return (c << i | c >>> 32 - i) + e
    }

    function ahh(t, e, n, r, o, i, a) {
      var c = t + (e ^ n ^ r) + (o >>> 0) + a;
      return (c << i | c >>> 32 - i) + e
    }

    function aii(t, e, n, r, o, i, a) {
      var c = t + (n ^ (e | ~r)) + (o >>> 0) + a;
      return (c << i | c >>> 32 - i) + e
    }

    function sign(t, n) {
      var e = n150;
      if (void 0 === t || null === t)
        throw new Error("Illegal argument " + t);
      var r = n150.wordsToBytes(a(t, n));
      return n && n.asBytes ? r : n && n.asString ? i.bytesToString(r) : e.bytesToHex(r)
    }

	result = sign(query);
`

func (*sogouTranslator) calSign(from, to, text string) (string, error) {
	vm := otto.New()
	if err := vm.Set("query", from+to+text+"109984457"); err != nil {
		return "", fmt.Errorf("vm.Set query error: %v", err)
	}
	value, err := vm.Run(signJS)
	if err != nil {
		return "", fmt.Errorf("vm.Run error: %v", err)
	}
	result, err := value.ToString()
	if err != nil {
		return "", fmt.Errorf("vlue.ToString error: %v", err)
	}
	return result, nil
}
