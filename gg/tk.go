package gg

import (
	"fmt"
	"regexp"

	"github.com/robertkrimen/otto"
)

const (
	tkJS = `
  var Tq = function(a) {
	        return function() {
	            return a
	        }
	    },
	    Uq = function(a, b) {
	        for (var c = 0; c < b.length - 2; c += 3) {
	            var d = b.charAt(c + 2);
	            d = "a" <= d ? d.charCodeAt(0) - 87 : Number(d);
	            d = "+" == b.charAt(c + 1) ? a >>> d : a << d;
	            a = "+" == b.charAt(c) ? a + d & 4294967295 : a ^ d
	        }
	        return a
	    },
	    Wq = function(a) {
	        var b = Vq;
	        var d = Tq(String.fromCharCode(116));
	        c = Tq(String.fromCharCode(107));
	        d = [d(), d()];
	        d[1] = c();
	        d = b.split(".");
	        b = Number(d[0]) || 0;
	        for (var e = [], f = 0, g = 0; g < a.length; g++) {
	            var l = a.charCodeAt(g);
	            128 > l ? e[f++] = l : (2048 > l ? e[f++] = l >> 6 | 192 : (55296 == (l & 64512) && g + 1 < a.length && 56320 == (a.charCodeAt(g + 1) & 64512) ? (l = 65536 + ((l & 1023) << 10) + (a.charCodeAt(++g) & 1023),
	                        e[f++] = l >> 18 | 240,
	                        e[f++] = l >> 12 & 63 | 128) : e[f++] = l >> 12 | 224,
	                    e[f++] = l >> 6 & 63 | 128),
	                e[f++] = l & 63 | 128)
	        }
	        a = b;
	        for (f = 0; f < e.length; f++)
	            a += e[f],
	            a = Uq(a, "+-a^+6");
	        a = Uq(a, "+-3^+b+-f");
	        a ^= Number(d[1]) || 0;
	        0 > a && (a = (a & 2147483647) + 2147483648);
	        a %= 1E6;
	        return a.toString() + "." + (a ^ b)
	    };
	result = Wq(query)
	`
)

func calTK(vq, query string) (string, error) {
	vm := otto.New()
	if err := vm.Set("Vq", vq); err != nil {
		return "", fmt.Errorf("vm.Set Vq error: %v", err)
	}
	if err := vm.Set("query", query); err != nil {
		return "", fmt.Errorf("vm.Set query error: %v", err)
	}
	value, err := vm.Run(tkJS)
	if err != nil {
		return "", fmt.Errorf("vm.Run error: %v", err)
	}
	return value.String(), nil
}

func getVq(dataStr string) (string, error) {
	//dataStr = `LOW_CONFIDENCE_THRESHOLD=-1;MAX_ALTERNATIVES_ROUNDTRIP_RESULTS=1;TKK=eval('((function(){var a\x3d1966732470;var b\x3d1714107181;return 423123+\x27.\x27+(a+b)})())');VERSION_LABEL = 'twsfe_w_20180402_RC00';`
	vqResult := regexp.MustCompile(`TKK=(?P<vq>[\s\S]+);VERSION_LABEL`).FindStringSubmatch(dataStr)
	if len(vqResult) != 2 {
		return "", fmt.Errorf("cannot get vq")
	}
	vm := otto.New()
	value, err := vm.Run(vqResult[1])
	if err != nil {
		return "", fmt.Errorf("vm.Run error: %v", err)
	}
	return value.String(), nil
}
