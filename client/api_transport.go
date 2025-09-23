package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type apiTransport struct {
	http.RoundTripper

	baseUri *url.URL
	apiKey  string
	debug   io.Writer
}

func (a *apiTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	if nil != a.baseUri {
		request.URL = a.baseUri.ResolveReference(request.URL)
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("api-key", a.apiKey)
	request.Header.Set("user-agent", "libdns-client/1.0")

	if nil != a.debug {
		dump(request, httputil.DumpRequest, "c", a.debug)
	}

	response, err := a.RoundTripper.RoundTrip(request)

	if nil != a.debug && nil != response {
		dump(response, httputil.DumpResponse, "s", a.debug)
	}

	return response, err
}

func dump[T *http.Request | *http.Response](x T, d func(T, bool) ([]byte, error), p string, o io.Writer) {
	if out, err := d(x, true); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			_, _ = fmt.Fprintf(o, "[%s] %s\n", p, scanner.Text())
		}
	}
}
