package main

import (
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

type Request struct {
	Method     string `json:"method"`
	URLPattern string `json:"urlPattern"`
}

func (r Request) Match(url, method string) bool {
	return r.Method == method && Must(regexp.MatchString(r.URLPattern, url))
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
