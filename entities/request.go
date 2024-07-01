package entities

import (
	"net/http"
	"regexp"

	"github.com/wesley601/mockid/utils"
)

type QueryParam struct {
	EqualTo string `json:"equalTo"`
	Matches string `json:"matches"`
}

type Request struct {
	Method          string                `json:"method"`
	URLPath         string                `json:"urlPath"`
	URLPattern      string                `json:"urlPattern"`
	QueryParameters map[string]QueryParam `json:"queryParameters"`
}

func (r Request) GetPath() string {
	if r.URLPath != "" {
		return r.URLPath
	}

	return r.URLPattern
}

func (r Request) Match(res *http.Request) bool {
	if r.Method != res.Method {
		return false
	}
	if len(r.QueryParameters) != 0 {
		for k, v := range r.QueryParameters {
			param := res.URL.Query().Get(k)
			if v.EqualTo != "" {
				if v.EqualTo != param {
					return false
				}
			}
			if v.Matches != "" {
				if !utils.Must(regexp.MatchString(v.Matches, param)) {
					return false
				}
			}
		}
	}
	if r.URLPath != "" {
		if r.URLPath != res.URL.Path {
			return false
		}
	}
	if !utils.Must(regexp.MatchString(r.URLPattern, res.URL.Path)) {
		return false
	}
	return true
}
