package rank

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func httpGet(fullUrl string) ([]byte, error) {
	client := &http.Client{}

	proxyStr := os.Getenv("HTTP_PROXY")
	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		client = &http.Client{Transport: transport}
	}
	// avatar server insecure
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodGet, fullUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
