package client

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFsClient(t *testing.T) {
	c := &FsClient{fs: afero.NewMemMapFs(), dir: "/tmp"}
	content := []byte(`{"name": "foo"}`)
	id, err := c.Create(content)
	require.NoError(t, err, "create failed")
	got, err := c.Read(id)
	require.Equal(t, content, got, "read after creation")
	content = []byte(`{"name": "bar"}`)
	require.NoError(t, c.Update(id, content), "update failed")
	got, err = c.Read(id)
	require.Equal(t, content, got, "read after update")
	require.NoError(t, c.Delete(id), "delete failed")
	_, err = c.Read(id)
	require.Equal(t, ErrNotFound, err, "read non existent resource should return ErrNotFound")
}
