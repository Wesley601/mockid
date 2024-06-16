package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

type Mappings struct {
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

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
				if !Must(regexp.MatchString(v.Matches, param)) {
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
	if !Must(regexp.MatchString(r.URLPattern, res.URL.Path)) {
		return false
	}
	return true
}

type Response struct {
	Headers      map[string]string `json:"headers"`
	BodyFileName string            `json:"bodyFileName"`
	Status       int               `json:"status"`
}

func (r Response) GetBody() ([]byte, error) {
	fullPath := path.Join("./mocks/__files/", r.BodyFileName)
	response, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type RequestSaved struct {
	ID              int
	RequestedPath   string
	RequestedMethod string
	MatchedPath     string
	ResponseBody    string
	ResponseStatus  int
	CreatedAt       string
}

type Scanner interface {
	Scan(dest ...any) error
}

func (r *RequestSaved) Scan(row Scanner) error {
	if err := row.Scan(&r.ID, &r.RequestedPath, &r.RequestedMethod, &r.MatchedPath, &r.ResponseBody, &r.ResponseStatus); err != nil {
		return fmt.Errorf("failed to scan a request row: %w", err)
	}
	return nil
}
