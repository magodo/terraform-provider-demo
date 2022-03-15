package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type testHandler struct {
	i   uint64
	buf map[string]map[string]interface{}
}

func (h *testHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("reading request: %v", err)))
			return
		}
		r.Body.Close()
		m := map[string]interface{}{}
		if err := json.Unmarshal(b, &m); err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unmarshal request: %v", err)))
			return
		}
		id := atomic.AddUint64(&h.i, 1)
		m["id"] = id
		url := joinPath(*r.URL, strconv.FormatUint(id, 10))
		h.buf[url.Path] = m
		b, err = json.Marshal(m)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("marshal response: %v", err)))
			return
		}
		w.WriteHeader(201)
		w.Write(b)
		return
	case http.MethodGet:
		m, ok := h.buf[r.URL.Path]
		if !ok {
			w.WriteHeader(404)
			return
		}
		b, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("marshal response: %v", err)))
			return
		}
		w.Write(b)
		return
	case http.MethodPut:
		om, ok := h.buf[r.URL.Path]
		if !ok {
			w.WriteHeader(404)
			return
		}
		id := om["id"]
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("reading request: %v", err)))
			return
		}
		r.Body.Close()
		m := map[string]interface{}{}
		if err := json.Unmarshal(b, &m); err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("unmarshal request: %v", err)))
			return
		}
		m["id"] = id
		h.buf[r.URL.Path] = m
		w.Write(b)
		return
	case http.MethodDelete:
		if _, ok := h.buf[r.URL.Path]; !ok {
			w.WriteHeader(404)
			return
		}
		delete(h.buf, r.URL.Path)
		w.WriteHeader(200)
		return
	}
}

func TestClientJSONServer(t *testing.T) {
	h := testHandler{
		buf: map[string]map[string]interface{}{},
	}
	ts := httptest.NewServer(http.HandlerFunc(h.Handle))
	defer ts.Close()
	c, _ := NewJSONServerClient(ts.URL + "/posts")

	id, err := c.Create([]byte(`{"name": "foo"}`))
	require.NoError(t, err, "create failed")
	got, err := c.Read(id)
	require.JSONEq(t, `{"id": 1, "name": "foo"}`, string(got), "read after creation")
	require.NoError(t, c.Update(id, []byte(`{"name": "bar"}`)), "update failed")
	got, err = c.Read(id)
	require.JSONEq(t, `{"id": 1, "name": "bar"}`, string(got), "read after update")
	require.NoError(t, c.Delete(id), "delete failed")
	_, err = c.Read(id)
	require.Error(t, err, "read after delete")
}
