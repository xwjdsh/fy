```
            ____
           / __/_  __
          / /_/ / / /
         / __/ /_/ /
        /_/  \__, /
            /____/
```
fy is a command-line tool for translation.

![](https://raw.githubusercontent.com/xwjdsh/fy/master/screenshot/fy.gif)

## Installation
### Homebrew
```
brew tap xwjdsh/tap
brew install xwjdsh/tap/fy
```
### Go
```
go get -u github.com/xwjdsh/cmd/fy
```
### Docker
```
alias fy=docker run -t --rm wendellsun/fy
```
### Manual
Download it from [releases](https://github.com/xwjdsh/fy/releases), and extact it to your `PATH` directory.

## Usage
```
Usage of ./fy:
  -d    Debug mode, if an error occurs in the translation, the error message is displayed
  -e string
        Select translators except these, comma separated. eg 'bd,gg', it can also be set by the 'FY_EXCEPT' environment variable
  -o string
        Select only the translators, comma separated. eg 'bd,gg', it can also be set by the 'FY_ONLY' environment variable
  -s    Display translators information
```
## Licence
[MIT License](https://github.com/xwjdsh/fy/blob/master/LICENSE)
