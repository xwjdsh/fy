package bd

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/xwjdsh/fy"
)

type baidu struct{}

func init() {
	fy.Translators = append(fy.Translators, &baidu{})
}

func (b *baidu) Translate(req *fy.Request) (resp *fy.Response) {
	resp = &fy.Response{}
	var from, to string
	if req.IsChinese {
		from, to = "zh", "en"
	} else {
		from, to = "en", "zh"
	}
	param := url.Values{
		"from":              {from},
		"to":                {to},
		"query":             {req.Text},
		"simple_means_flag": {"3"},
		"sign":              {""},
		"token":             {""},
	}
	r, err := http.PostForm("http://fanyi.baidu.com/v2transapi", param)
	if err != nil {
		resp.Err = errors.Wrap(err, "http.PostForm error")
		return
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Err = errors.Wrap(err, "ioutil.ReadAll error")
		return
	}
	log.Println(string(body))
	return
}
