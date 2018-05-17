```
            ____
           / __/_  __
          / /_/ / / /
         / __/ /_/ /
        /_/  \__, /
            /____/
```
[![Build Status](https://travis-ci.org/xwjdsh/fy.svg?branch=master)](https://travis-ci.org/xwjdsh/fy)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/fy)](https://goreportcard.com/report/github.com/xwjdsh/fy)
[![GoDoc](https://godoc.org/github.com/xwjdsh/fy?status.svg)](https://godoc.org/github.com/xwjdsh/fy)
[![](https://images.microbadger.com/badges/image/wendellsun/fy.svg)](https://microbadger.com/images/wendellsun/fy)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/fy/blob/master/LICENSE)


中文 | [English](https://github.com/xwjdsh/fy/blob/master/README_EN.md)

`fy`是一个命令行下的翻译工具。

![](https://raw.githubusercontent.com/xwjdsh/fy/master/screenshot/fy.gif)
## 安装
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
从 [releases](https://github.com/xwjdsh/fy/releases) 下载可执行文件并将其放到 PATH 环境变量对应的路径中。

## 使用
```
Usage of fy:
  -d    调试模式，如果翻译过程出现错误，会将错误信息作为翻译结果展示
  -e string
        选择除了指定以外的翻译者, 逗号分隔, 例如 'bd,gg', 也可以通过 'FY_EXCEPT' 环境变量来配置
  -f string
        翻译文件的路径
  -o string
        选择指定的翻译者, 逗号分隔, 例如 'bd,gg', 也可以通过 'FY_ONLY' 环境变量来配置
  -s    显示支持的翻译者的信息
  -t    指定翻译的目标语言
```

### 语言

| 简写 | 对应语言 | 
| - | :-: | 
| zh-CN | 中文 | 
| en | 英语 | 
| ru | 俄语 | 
| ja | 日语 | 
| de | 德语 | 
| fr | 法语 | 
| ko | 韩语 | 
| es | 西班牙语 | 

### 示例
```shell
# 显示支持的翻译者的信息
fy -s

# 不输入参数，则会访问剪贴板内容
fy

# 普通方式
fy test

# 调试模式
fy -d test

# 指定目标语言
fy -t ja 测试翻译为日语

# 指定中文的翻译语言，默认为英语
FY_CN_TO=ko fy 翻译为韩语

# 指定非中文的翻译语言，默认为简体中文
FY_NOT_CN_TO=en fy 중국어로 번역

# 翻译文件
cat `test.txt` | fy
fy < test.txt
fy -f test.txt

# 选择除了指定以外的翻译者
FY_EXCEPT='bd,sg' fy test
fy -e 'bd,sg' test

# 选择指定的翻译者
FY_ONLY='gg,qq' fy test
fy -o 'gg,qq' test
```

## 协议
[MIT License](https://github.com/xwjdsh/fy/blob/master/LICENSE)
