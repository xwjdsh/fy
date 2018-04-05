package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/xwjdsh/fy"
	_ "github.com/xwjdsh/fy/bd"
	_ "github.com/xwjdsh/fy/sg"
)

var (
	version string
	isDebug = flag.Bool("d", false, "debug")
)

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
	flag.Parse()
	args := os.Args
	if len(args) == 1 {
		fmt.Printf(fy.Logo, version)
		return
	}
	text := strings.Join(args[1:], " ")
	isChinese := fy.IsChinese(text)

	for _, t := range fy.Translators {
		go fy.Handle(t, &fy.Request{
			IsChinese: isChinese,
			Text:      text,
		})
	}

	for range fy.Translators {
		resp := <-fy.ResponseChan
		log.Println(resp)
	}
}
