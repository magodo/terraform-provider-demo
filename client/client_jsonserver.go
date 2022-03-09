package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type JSONServerClient struct {
	baseURL url.URL
}

func NewJSONServerClient(endpoint string) (Client, error) {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	return &JSONServerClient{
		baseURL: *baseURL,
	}, nil
}

func statuscodeOK(code int, okCodes ...int) bool {
	for _, okcode := range okCodes {
		if okcode == code {
			return true
		}
	}
	return false
}

func (j *JSONServerClient) Create(b []byte) (string, error) {
	resp, err := http.Post(j.baseURL.String(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("post: %w", err)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if !statuscodeOK(resp.StatusCode, http.StatusOK, http.StatusCreated) {
		return "", fmt.Errorf("unexpected status code: %d. Message: %s", resp.StatusCode, string(content))
	}
	payload := map[string]interface{}{}
	if err := json.Unmarshal(content, &payload); err != nil {
		return "", err
	}
	return strconv.Itoa(int(payload["id"].(float64))), nil
}

func (j *JSONServerClient) Read(id string) ([]byte, error) {
	url, err := j.baseURL.Parse(id)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if !statuscodeOK(resp.StatusCode, http.StatusOK) {
		return nil, fmt.Errorf("unexpected status code: %d. Message: %s", resp.StatusCode, string(content))
	}
	return content, nil
}

func (j *JSONServerClient) Update(id string, b []byte) error {
	url, err := j.baseURL.Parse(id)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !statuscodeOK(resp.StatusCode, http.StatusOK) {
		return fmt.Errorf("unexpected status code: %d. Message: %s", resp.StatusCode, string(content))
	}
	return nil
}

func (j *JSONServerClient) Delete(id string) error {
	url, err := j.baseURL.Parse(id)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !statuscodeOK(resp.StatusCode, http.StatusOK) {
		return fmt.Errorf("unexpected status code: %d. Message: %s", resp.StatusCode, string(content))
	}
	return nil
}
