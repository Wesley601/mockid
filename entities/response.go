package entities

import (
	"os"
	"path"

	"github.com/wesley601/mockid/utils"
)

type Response struct {
	Headers      map[string]string `json:"headers"`
	BodyFileName string            `json:"bodyFileName"`
	Status       int               `json:"status"`
}

func (r Response) GetBody() ([]byte, error) {
	return GetBody(r.BodyFileName)
}

func GetBody(bodyFileName string) ([]byte, error) {
	fullPath := path.Join(utils.Must(os.Getwd()), "/mocks/__files/", bodyFileName)
	response, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return response, nil
}
