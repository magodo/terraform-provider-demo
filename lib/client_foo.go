package lib

import (
	"encoding/json"
	"github.com/getlantern/deepcopy"
	"io/ioutil"
	"os"
)

type ClientFoo struct{}
type ModelFoo struct {
	Id      *string     `json:"id"`
	Name    *string     `json:"name"`
	Age     *int        `json:"age"`
	Job     *string     `json:"job"`
	Aliases *[]string   `json:"aliases"`
	Contact *ContactFoo `json:"contact"`
	Addrs   *[]*Addr    `json:"addrs"`
}

type ContactFoo struct {
	Phone  *int    `json:"phone"`
	Github *string `json:"github"`
}

type Addr struct {
	Country *string   `json:"country"`
	City    *string   `json:"city"`
	Roads   *[]string `json:"roads"`
}

const storageFoo = "/tmp/resource_foo.json"

func (c *ClientFoo) CreateOrUpdate(req *ModelFoo) (*ModelFoo, error) {
	var resp ModelFoo
	if err := deepcopy.Copy(&resp, req); err != nil {
		return nil, err
	}
	resp.Id = req.Name

	// store in fs
	b, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(storageFoo, b, 0644); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *ClientFoo) Get(id string) (*ModelFoo, error) {
	// fetch from fs
	b, err := ioutil.ReadFile(storageFoo)
	if err != nil {
		return nil, err
	}
	var resp ModelFoo
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *ClientFoo) Delete(id string) error {
	return os.Remove(storageFoo)
}
