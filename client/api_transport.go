package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type apiTransport struct {
	http.RoundTripper
	ApiClientConfig
}

func (a *apiTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	if uri := a.GetBaseUri(); uri != nil {
		request.URL = uri.ResolveReference(request.URL)
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("api-key", a.GetApiKey())
	request.Header.Set("user-agent", "libdns-client/1.0")

	if writer := a.GetDebug(); writer != nil {
		dump(request, httputil.DumpRequest, "c", writer)
	}

	response, err := a.RoundTripper.RoundTrip(request)

	if writer := a.GetDebug(); writer != nil && nil != response {
		dump(response, httputil.DumpResponse, "s", writer)
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
