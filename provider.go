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
	ApiKey string `json:"api_key,omitempty"`
	// Debug when true it will dump the http.Client request/response to os.Stdout
	// for now you can change that by calling `SetDebug(w io.Writer)` direct
	Debug bool `json:"debug"`
	// BaseUri used for the api calls and will default to https://mijn.host/api/v2/
	BaseUri string `json:"base_uri,omitempty"`

	client *client.ApiClient
	mutex  sync.RWMutex
}

func (p *Provider) getClient() *client.ApiClient {

	if nil == p.client {
		p.client = client.NewApiClient("", nil)
		p.ReloadConfig()
	}

	return p.client
}

func (p *Provider) ReloadConfig() {

	if p.Debug {
		p.client.SetDebug(os.Stdout)
	} else {
		p.client.SetDebug(nil)
	}

	if "" != p.ApiKey {
		p.client.SetApiKey(p.ApiKey)
	}

	if "" != p.BaseUri {
		_ = p.client.SetBaseUrl(p.BaseUri)
	}

	p.client.CloseIdleConnections()
	p.ApiKey = ""
	p.BaseUri = ""
}

func (p *Provider) GetBaseUrl() *url.URL {
	return p.getClient().GetBaseUrl()
}

func (p *Provider) GetApiKey() string {
	return p.getClient().GetApiKey()
}

func (p *Provider) SetDebug(writer io.Writer) {
	p.Debug = writer != nil
	p.getClient().SetDebug(writer)
}

func (p *Provider) IsDebug() bool {
	return nil != p.getClient().GetDebug()
}

func fqdn(name string) string {

	if name[len(name)-1] != '.' {
		return name + "."
	}

	return name
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
	_ libdns.ZoneLister     = (*Provider)(nil)
)
