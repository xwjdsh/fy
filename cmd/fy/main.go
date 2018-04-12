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
	only    = flag.String("o", "", "only")
	except  = flag.String("e", "", "except")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(os.Args) == 1 || len(args) == 0 {
		fmt.Printf(fy.Logo, version)
		return
	}
	text := strings.Join(args, " ")
	isChinese := fy.IsChinese(text)

	if *only == "" {
		*only = os.Getenv("FY_ONLY")
	}
	if *except == "" {
		*except = os.Getenv("FY_EXCEPT")
	}
	translators, err := fy.Filter(*only, *except)
	if err != nil {
		color.Red("✗ %v", err)
		return
	}

	req := &fy.Request{
		IsChinese: isChinese,
		Text:      text,
	}
	responseChan := make(chan *fy.Response)

	wrap := func(t fy.Translator) {
		responseChan <- t.Translate(req)
	}
	for _, t := range translators {
		go wrap(t)
	}

	fmt.Println()
	for range translators {
		resp := <-responseChan
		if resp.Err != nil {
			if !*isDebug {
				continue
			}
			resp.Result = resp.Err.Error()
		}
		color.Green("\t%s  [%s]\n\n", fy.CoffeeEmoji, resp.FullName)
		color.Magenta("\t\t%s\n\n", resp.Result)
	}
}
