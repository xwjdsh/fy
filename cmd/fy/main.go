package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/xwjdsh/fy"
	_ "github.com/xwjdsh/fy/bd"
	_ "github.com/xwjdsh/fy/gg"
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

	wrap := func(t fy.Translator) {
		responseChan <- t.Translate(req)
	}
	for _, t := range fy.TranslatorMap {
		go wrap(t)
	}

	fmt.Println()
	for range fy.TranslatorMap {
		resp := <-responseChan
		color.Green("\t%s  [%s]\n\n", fy.CoffeeEmoji, resp.FullName)
		color.Magenta("\t\t%s\n\n", resp.Result)
	}
}
