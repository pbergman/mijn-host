package client

import (
	"io"
	"net/url"
)

type ApiClientConfig interface {
	GetApiKey() string
	GetDebug() io.Writer
	GetBaseUri() *url.URL
}
