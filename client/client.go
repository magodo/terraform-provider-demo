package client

type Client interface {
	Create(b []byte) (id string, err error)
	Read(id string) ([]byte, error)
	Update(id string, b []byte) error
	Delete(id string) error
}
