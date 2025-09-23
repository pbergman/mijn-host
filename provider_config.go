package mijnhost

import (
	"os"

	"github.com/pbergman/mijnhost/client"
)

type config interface {
	getApiKey() string
	isDebug() bool
	getBaseUri() string
}

type configWrapper struct {
	p *Provider
}

func (p *configWrapper) getApiKey() string {
	return p.p.ApiKey
}

func (p *configWrapper) isDebug() bool {
	return p.p.Debug
}

func (p *configWrapper) getBaseUri() string {
	return p.p.BaseUri
}

type configApi struct {
	ApiKey  string `json:"api_key"`
	Debug   bool   `json:"debug"`
	BaseUri string `json:"base_uri"`
}

func (p *configApi) getApiKey() string {
	return p.ApiKey
}

func (p *configApi) isDebug() bool {
	return p.Debug
}

func (p *configApi) getBaseUri() string {
	return p.BaseUri
}

func reload(p *Provider, data config) {

	if nil == p.client {
		p.client = client.NewApiClient("", nil)
	} else {
		p.client.CloseIdleConnections()
	}

	if data.isDebug() {
		p.SetDebug(os.Stdout)
	} else {
		p.SetDebug(nil)
	}

	if key := data.getApiKey(); key != "" {
		p.client.SetApiKey(key)
	}

	if uri := data.getBaseUri(); uri != "" {
		_ = p.client.SetBaseUrl(uri)
	}

	p.ApiKey = ""
	p.BaseUri = ""
}
