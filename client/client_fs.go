package client

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-uuid"
	"github.com/spf13/afero"
)

type FsClient struct {
	fs  afero.Fs
	dir string
}

func NewFsClient(dir string) (Client, error) {
	return &FsClient{fs: afero.NewOsFs(), dir: dir}, os.MkdirAll(dir, 0755)
}

func (f *FsClient) Create(b []byte) (string, error) {
	// We should check duplication of the generated filename (i.e. the UUID) in the directory.
	// In fact we shall use the os.CreateTemp() instead. However, since we are also using the afero
	// to make the UT less dependent to the OS, and afero.Fs doesn't implemented the CreateTemp().
	// For the sake of simplicity, we choose to use the UUID here for now.
	id, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	file, err := f.fs.Create(filepath.Join(f.dir, id))
	if err != nil {
		return "", err
	}
	defer file.Close()
	return id, f.Update(id, b)
}

func (f *FsClient) Update(id string, b []byte) error {
	return afero.WriteFile(f.fs, filepath.Join(f.dir, id), b, 0666)
}

func (f *FsClient) Read(id string) ([]byte, error) {
	b, err := afero.ReadFile(f.fs, filepath.Join(f.dir, id))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	return b, err
}

func (f *FsClient) Delete(id string) error {
	return f.fs.Remove(filepath.Join(f.dir, id))
}
