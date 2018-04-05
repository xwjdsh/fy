package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xwjdsh/fy"
	_ "github.com/xwjdsh/fy/sogou"
)

var version string

func init() {
	if version == "" {
		if output, err := exec.Command("git describe --tags").Output(); err == nil {
			version = string(output)
		} else {
			version = "unknown"
		}
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf(fy.Logo, version)
		return
	}
	text := strings.Join(args[1:], " ")
	isChinese := fy.IsChinese(text)
	for _, t := range fy.Translators {
		go t.Translate(isChinese, text)
	}
	for range fy.Translators {
	}
}
