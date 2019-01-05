package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HTTPClient struct {
	URL     string
	Path    string
	Method  string
	Headers map[string]string
	Param   url.Values
	JSON    interface{}
	IsJSON  bool
	Timeout time.Duration
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Headers: map[string]string{},
	}
}

// Get is specify for HTTP POST Method with JSON body.
func (r *HTTPClient) JSONPost(u *url.URL) (*http.Request, error) {
	u.Path += r.Path
	link := u.String()
	body, err := json.Marshal(r.JSON)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method, link, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r.Headers["Content-Type"] = "application/json"

	return req, nil
}

// Get is specify for HTTP POST Method.
func (r *HTTPClient) Post(u *url.URL) (*http.Request, error) {
	if r.IsJSON {
		return r.JSONPost(u)
	}
	u.Path += r.Path
	link := u.String()
	form := strings.NewReader(r.Param.Encode())
	req, err := http.NewRequest(r.Method, link, form)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	return req, nil
}

// Get is specify for HTTP GET Method.
func (r *HTTPClient) Get(u *url.URL) (*http.Request, error) {
	u.RawQuery = r.Param.Encode()
	u.Path += r.Path
	link := u.String()
	req, err := http.NewRequest(r.Method, link, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return req, nil
}

// DoReq is Last Point to call api
func (r *HTTPClient) DoReq() (*[]byte, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var req *http.Request
	switch r.Method {
	case "GET":
		req, err = r.Get(u)
	case "POST":
		req, err = r.Post(u)
	}

	if req == nil {
		return nil, errors.New("Failed create new request.")
	}

	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	req.Close = true

	if r.Timeout < 1 {
		r.Timeout = time.Duration(20 * time.Second)
	}
	hc := &http.Client{
		Timeout: r.Timeout,
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.Println(resp, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("%+v\n", resp)
		log.Printf("%+v\n", req)
		log.Println(u.String())
		return nil, fmt.Errorf("Status Code = %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &contents, nil
}
