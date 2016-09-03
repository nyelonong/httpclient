package httpclient

import (
	"log"
	"net/url"
	"testing"
	"time"
)

const (
	URL string = "https://httpbin.org"
)

func TestRun(t *testing.T) {
	param := url.Values{}
	param.Add("show_env", "1")

	headers := make(map[string]string)
	headers["Foo"] = "Bar"

	hc := NewHTTPClient()
	hc.URL = URL
	hc.Path = "get"
	hc.Method = "GET"
	hc.Param = param
	hc.Headers = headers
	hc.Timeout = 5 * time.Second

	body, err := hc.DoReq()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(*body))
}
