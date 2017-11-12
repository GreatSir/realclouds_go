package utils

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//HTTPLib *
type HTTPLib struct {
	URL    string
	Client *http.Client
}

//NewHTTPLib *
func NewHTTPLib(url string) *HTTPLib {
	url = strings.TrimSpace(url)
	url = strings.TrimRight(url, "/")

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}

	httpLib := &HTTPLib{
		URL:    url,
		Client: client,
	}

	return httpLib
}

//POST *
func (h *HTTPLib) POST(path string, body io.Reader) ([]byte, error) {
	return h.Do("POST", path, body)
}

//GET *
func (h *HTTPLib) GET(path string, body io.Reader) ([]byte, error) {
	return h.Do("GET", path, body)
}

//Do *
func (h *HTTPLib) Do(method, path string, body io.Reader) ([]byte, error) {
	path = strings.TrimLeft(path, "/")
	reqURL := h.URL + "/" + path

	if strings.Contains(h.URL, "?") {
		reqURL = h.URL + path
	}

	log.Printf("Request uri:%v\n", reqURL)

	req, err := http.NewRequest(method, reqURL, body)
	if nil != err {
		return nil, err
	}

	req.Header.Set("User-Agent", "AyerDudu http client v1.0")
	req.Header.Set("Auth", "Eric Shi / shibingli@yeah.net")

	resp, err := h.Client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}
