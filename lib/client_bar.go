package lib

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	"github.com/getlantern/deepcopy"
)

type ClientBar struct{}
type ModelBar struct {
	Id     *string `json:"id"`
	Name   *string `json:"name"`
	Github *string `json:"github"`
	Phone  *int    `json:"phone"`
}

const storageBar = "/tmp/resource_bar.json"

func (c *ClientBar) CreateOrUpdate(req *ModelBar) (*ModelBar, error) {
	var resp ModelBar
	if err := deepcopy.Copy(&resp, req); err != nil {
		return nil, err
	}
	id := strconv.Itoa(rand.Int())
	resp.Id = &id

	// store in fs
	b, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(storageBar, b, 0644); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *ClientBar) Get(id string) (*ModelBar, error) {
	// fetch from fs
	b, err := ioutil.ReadFile(storageBar)
	if err != nil {
		return nil, err
	}
	var resp ModelBar
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *ClientBar) Delete(id string) error {
	return os.Remove(storageBar)
}
