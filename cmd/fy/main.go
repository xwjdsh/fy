package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/xwjdsh/fy"
	_ "github.com/xwjdsh/fy/bd"
	_ "github.com/xwjdsh/fy/by"
	_ "github.com/xwjdsh/fy/gg"
	_ "github.com/xwjdsh/fy/qq"
	_ "github.com/xwjdsh/fy/sg"
	_ "github.com/xwjdsh/fy/yd"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	version    = "unknown"
	isDebug    = flag.Bool("d", false, "Debug mode, if an error occurs in the translation, the error message is displayed")
	sources    = flag.Bool("s", false, "Display translators information")
	filePath   = flag.String("f", "", "file path")
	only       = flag.String("o", "", "Select only the translators, comma separated. eg 'bd,gg', it can also be set by the 'FY_ONLY' environment variable")
	except     = flag.String("e", "", "Select translators except these, comma separated. eg 'bd,gg', it can also be set by the 'FY_EXCEPT' environment variable")
	targetLang = flag.String("t", "", "The target language of translation")
)

func main() {
	flag.Parse()
	if *sources {
		printSources()
		return
	}

	text, err := getText()
	if err != nil {
		color.Red("%s %v", fy.IconBad, err)
		os.Exit(1)
	}
	if text == "" {
		color.Green(fy.Logo)
		fmt.Printf(fy.Desc, version)
		return
	}
	isChinese := fy.IsChinese(text)
	if *targetLang == "" {
		*targetLang = getTargetLang(isChinese)
	}

	translators, err := getTranslators()
	if err != nil {
		color.Red("%s %v", fy.IconBad, err)
		os.Exit(1)
	}

	req := fy.Request{
		TargetLang: *targetLang,
		Text:       text,
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

func getText() (string, error) {
	var text string
	if *filePath != "" {
		data, err := ioutil.ReadFile(*filePath)
		if err != nil {
			return "", err
		}
		text = string(data)
	} else if !terminal.IsTerminal(0) {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		text = string(data)
	} else {
		args := flag.Args()
		if len(os.Args) == 1 || len(args) == 0 {
			return "", nil
		}
		text = strings.Join(args, " ")
	}
	return text, nil
}

func getTargetLang(isChinese bool) string {
	target := os.Getenv("FY_TO")
	if target != "" {
		return target
	}
	if isChinese {
		target = os.Getenv("FY_CN_TO")
	} else {
		target = os.Getenv("FY_NOT_CN_TO")
	}
	if target == "" {
		if isChinese {
			target = "en"
		} else {
			target = "zh-CN"
		}
	}
	return target
}

func getTranslators() ([]fy.Translator, error) {
	if *only == "" {
		*only = os.Getenv("FY_ONLY")
	}
	if *except == "" {
		*except = os.Getenv("FY_EXCEPT")
	}
	return fy.Filter(*only, *except)
}

func printSources() {
	translators, _ := fy.Filter("", "")
	fmt.Println()
	for _, t := range translators {
		fy.PrintSource(t.Desc())
	}
	fmt.Println()
}
