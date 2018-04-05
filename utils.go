package fy

import "unicode"

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
