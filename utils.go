package fy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"unicode"
)

var (
	Logo = `
	    ____     
	   / __/_  __
	  / /_/ / / /
	 / __/ /_/ / 
	/_/  \__, /  
	    /____/   

 version: %s
homepage: https://github.com/xwjdsh/fy

`
)

func IsChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func ReadResp(resp *http.Response, err error) (*http.Response, []byte, error) {
	if err != nil {
		return nil, nil, fmt.Errorf("http request error: %v\n", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("ioutil.ReadAll error: %v\n", err)
	}
	return resp, respBody, nil
}

func NotReadResp(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, fmt.Errorf("http response error: %v\n", err)
	}
	defer resp.Body.Close()
	return resp, nil
}

func SendRequest(method, urlStr string, body io.Reader, f func(*http.Request) error) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, nil, fmt.Errorf("http.NewRequest error: %v\n", err)
	}
	client := &http.Client{}
	if f != nil {
		if err := f(req); err != nil {
			return nil, nil, fmt.Errorf("f error: %v\n", err)
		}
	}
	return ReadResp(client.Do(req))
}

func AddCookies(req *http.Request, cookies []*http.Cookie) *http.Request {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return req
}
