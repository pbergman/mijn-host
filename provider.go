package mijnhost

import (
	"io"
	"net/url"
	"os"
	"sync"

	"github.com/libdns/libdns"
	"github.com/pbergman/mijnhost/client"
)

type Provider struct {
	// ApiKey used for authenticating the mijn.host api see:
	// https://mijn.host/api/doc/doc-343216#obtaining-your-api-key
	ApiKey string `json:"api_key"`
	// Debug when true it will dump the http.Client request/response to os.Stdout
	// or you can change that by setting `DebugOut`
	Debug    bool      `json:"debug"`
	DebugOut io.Writer `json:"-"`
	// BaseUri used for the api calls and will default to https://mijn.host/api/v2/
	BaseUri *ApiBaseUri `json:"base_uri"`

	client *client.ApiClient
	mutex  sync.RWMutex
}

func (p *Provider) getClient() *client.ApiClient {
	if nil == p.client {

		if p.BaseUri == nil {
			p.BaseUri = DefaultApiBaseUri()
		}

		p.client = client.NewApiClient(p)
	}

	return p.client
}

func (o *Provider) GetApiKey() string {
	return o.ApiKey
}

func (o *Provider) GetDebug() io.Writer {
	if o.Debug {
		if nil == o.DebugOut {
			return os.Stdout
		}
		return o.DebugOut
	}
	return nil
}

func (o *Provider) GetBaseUri() *url.URL {

	if nil == o.BaseUri {
		return nil
	}

	return (*url.URL)(o.BaseUri)
}

func fqdn(name string) string {

	if name[len(name)-1] != '.' {
		return name + "."
	}

	return name
}

// Interface guards
var (
	_ client.ApiClientConfig = (*Provider)(nil)
	_ libdns.RecordGetter    = (*Provider)(nil)
	_ libdns.RecordAppender  = (*Provider)(nil)
	_ libdns.RecordSetter    = (*Provider)(nil)
	_ libdns.RecordDeleter   = (*Provider)(nil)
	_ libdns.ZoneLister      = (*Provider)(nil)
)
