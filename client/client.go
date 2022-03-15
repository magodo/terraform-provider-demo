package client

import "errors"

// ErrNotFound is expected to be returned for `Read` when the resource with the specified id doesn't exist.
var ErrNotFound = errors.New("resource not found")

type Client interface {
	Create(b []byte) (id string, err error)
	Read(id string) ([]byte, error)
	Update(id string, b []byte) error
	Delete(id string) error
}
