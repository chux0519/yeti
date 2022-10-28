package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func HttpPost(url string, data []byte, timeout time.Duration) ([]byte, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(body), "error") {
		return nil, errors.New(string(body))
	}

	return body, err
}
