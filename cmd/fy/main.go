package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/xwjdsh/fy"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	iconBad = "✗"

	logo = `
	    ____     
	   / __/_  __
	  / /_/ / / /
	 / __/ /_/ / 
	/_/  \__, /  
	    /____/   

https://github.com/xwjdsh/fy
`
	coffeeEmoji = "\u2615\ufe0f"
)

var (
	isDebug    = flag.Bool("d", false, "Debug mode, if an error occurs during translation, the error message will be displayed as the translation result")
	filePath   = flag.String("f", "", "File path")
	targetLang = flag.String("t", "", "The target language of the translation")

	translator = flag.String("translator", "", "Restrict the translators used, comma separated. eg 'baidu,google'")
	timeout    = flag.Duration("timeout", 5*time.Second, "The timeout for each translator")
)

func main() {
	flag.Parse()

	text, err := getText()
	if err != nil {
		color.Red("%s %v", iconBad, err)
		os.Exit(1)
	}
	if text == "" {
		color.Green(logo)
		return
	}
	isChinese := fy.IsChinese(text)
	if *targetLang == "" {
		*targetLang = getTargetLang(isChinese)
	}

	req := &fy.Request{
		ToLang: *targetLang,
		Text:   text,
	}

	var translators []string
	if *translator != "" {
		translators = strings.Split(*translator, ",")
	}

	ch := fy.AsyncTranslate(*timeout, req, translators...)
	for resp := range ch {
		if resp.Err != nil {
			if !*isDebug {
				continue
			}
			resp.Result = resp.Err.Error()
		}
		color.Green("\t%s  [%s]\n\n", coffeeEmoji, resp.Name)
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
			return clipboard.ReadAll()
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
