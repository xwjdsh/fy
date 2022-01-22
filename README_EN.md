```
            ____
           / __/_  __
          / /_/ / / /
         / __/ /_/ /
        /_/  \__, /
            /____/
```
[中文](https://github.com/xwjdsh/fy/blob/master/README.md) | English

[![Release](https://img.shields.io/github/release/xwjdsh/fy.svg?style=flat-square)](https://github.com/xwjdsh/fy/releases/latest)
[![Build Status](https://travis-ci.org/xwjdsh/fy.svg?branch=master)](https://travis-ci.org/xwjdsh/fy)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/fy)](https://goreportcard.com/report/github.com/xwjdsh/fy)
[![GoDoc](https://godoc.org/github.com/xwjdsh/fy?status.svg)](https://godoc.org/github.com/xwjdsh/fy)
[![](https://images.microbadger.com/badges/image/wendellsun/fy.svg)](https://microbadger.com/images/wendellsun/fy)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/fy/blob/master/LICENSE)

`fy` is a command-line tool for translation.

![](https://raw.githubusercontent.com/xwjdsh/fy/master/screenshot/fy.gif)

## Statement
This project is written for learning purposes only, the act of acquiring and sharing may be suspected of violating rights. please understand the situation, if your rights and interests are infringed, please contact me and I will delete it in time.

## Installation
### Homebrew
```
brew tap xwjdsh/tap
brew install xwjdsh/tap/fy
```
### Go
```
go get -u github.com/xwjdsh/fy/cmd/fy
```
### Docker
```
alias fy='docker run -t --rm wendellsun/fy'
```
### Manual
Download it from [releases](https://github.com/xwjdsh/fy/releases), and extact it to your `PATH` directory.

## Usage
```
Usage of ./fy:
  -d    Debug mode, if an error occurs during translation, the error message will be displayed as the translation result
  -f string
        File path
  -s    Display translator sources
  -t string
        The target language of the translation
  -timeout duration
        The timeout for each translator (default 5s)
  -translator string
        Restrict the translators used, comma separated. eg 'baidu,google'
```

### Language Mapping

| Shorthand | Language | 
| - | :-: | 
| zh-CN | Chinese | 
| en | English | 
| ru | Russian | 
| ja | Japanese | 
| de | German | 
| fr | French | 
| ko | Korean | 
| es | Spanish | 

### Example
```shell
# display supported translators
fy -s

# if there are no params, the clipboard will be accessed
fy

# simplest
fy test

# debug mode
fy -d test

# specify the language
fy -t ja 测试翻译为日语

# target language for Chinese，default is English
FY_CN_TO=ko fy 翻译为韩语

# target language for Non-Chinese，default is simplified Chinese
FY_NOT_CN_TO=en fy 중국어로 번역

# for file
cat `test.txt` | fy
fy < test.txt
fy -f test.txt

# select only the translators
fy -translator 'baidu,google' test
```

## Licence
[MIT License](https://github.com/xwjdsh/fy/blob/master/LICENSE)
