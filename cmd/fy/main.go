package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xwjdsh/fy"
	_ "github.com/xwjdsh/fy/bd"
	_ "github.com/xwjdsh/fy/sg"
	_ "github.com/xwjdsh/fy/yd"
)

var (
	version = "unknown"
	isDebug = flag.Bool("d", false, "debug")
)

func main() {
	flag.Parse()
	args := os.Args
	if len(args) == 1 {
		fmt.Printf(fy.Logo, version)
		return
	}
	text := strings.Join(args[1:], " ")
	isChinese := fy.IsChinese(text)

	req := &fy.Request{
		IsChinese: isChinese,
		Text:      text,
	}

	responseChan := make(chan *fy.Response)

	wrap := func(t fy.Translator, r *fy.Request) {
		responseChan <- t.Translate(r)
	}
	for _, t := range fy.Translators {
		go wrap(t, req)
	}

	for range fy.Translators {
		resp := <-responseChan
		log.Println(resp.Name, resp.Result)
	}
}
