package mijnhost

import (
	"encoding/json"
	"io"
	"net/url"
	"sync"

	"github.com/libdns/libdns"
	"github.com/pbergman/mijnhost/client"
)

type Provider struct {
	// ApiKey used for authenticating the mijn.host api see:
	// https://mijn.host/api/doc/doc-343216#obtaining-your-api-key
	ApiKey string
	// Debug when true it will dump the http.Client request/response to os.Stdout
	// for now you can change that by calling `SetDebug(w io.Writer)` direct
	Debug bool
	// BaseUri used for the api calls and will default to https://mijn.host/api/v2/
	BaseUri string

	client *client.ApiClient
	mutex  sync.RWMutex
}

func (p *Provider) MarshalJSON() ([]byte, error) {
	if p.client != nil {

		var object = configApi{
			ApiKey:  p.client.GetApiKey(),
			Debug:   p.Debug,
			BaseUri: p.client.GetBaseUrl().String(),
		}

		return json.Marshal(object)
	} else {
		return json.Marshal(p)
	}
}

func (p *Provider) UnmarshalJSON(data []byte) error {
	var object configApi

	if err := json.Unmarshal(data, &object); err != nil {
		return err
	}

	reload(p, &object)

	return nil
}

func (p *Provider) getClient() *client.ApiClient {
	if nil == p.client {
		p.ReloadConfig()
	}

	return p.client
}

func (p *Provider) ReloadConfig() {
	reload(p, &configWrapper{p})
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
