package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/getlantern/deepcopy"
)

type ClientBar struct{}
type ModelBar struct {
	Id    *string `json:"id"`
	Name  *string `json:"name"`
	Job   *string `json:"job"`
	Phone *int    `json:"phone"`
}

func storageBar(name string) string {
	return fmt.Sprintf("/tmp/resource_bar_%s.json", name)
}

func (c *ClientBar) CreateOrUpdate(req *ModelBar) (*ModelBar, error) {
	var resp ModelBar
	if err := deepcopy.Copy(&resp, req); err != nil {
		return nil, err
	}
	resp.Id = req.Name

	// store in fs
	b, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(storageBar(*req.Name), b, 0644); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *ClientBar) Get(id string) (*ModelBar, error) {
	// fetch from fs
	b, err := ioutil.ReadFile(storageBar(id))
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
	return os.Remove(storageBar(id))
}
