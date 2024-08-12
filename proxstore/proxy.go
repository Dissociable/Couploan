package proxstore

import (
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/phuslu/shardmap"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type Proxy[C any] struct {
	Protocol          Protocol
	Host              string
	Port              uint16
	Username          string
	Password          string
	Rotating          bool
	httpClient        *shardmap.Map[string, C]
	httpClientCreator CreateHttpClientCreator[C]
}

// NewProxy creates a new proxy
func NewProxy[C any](host string, port uint16, protocol Protocol) *Proxy[C] {
	return &Proxy[C]{
		Protocol: protocol,
		Host:     host,
		Port:     port,
	}
}

// NewProxyWithCredential creates a new proxy with credential
func NewProxyWithCredential[C any](
	host string, port uint16, protocol Protocol, username string, password string,
) *Proxy[C] {
	return &Proxy[C]{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SetHttpClientCreator sets the http client creator
func (p *Proxy[C]) SetHttpClientCreator(creator CreateHttpClientCreator[C]) *Proxy[C] {
	p.httpClientCreator = creator
	return p
}

// String returns the string representation of the proxy, e.g., http://username:password@host:port
func (p *Proxy[C]) String() string {
	if p.Protocol == ProtocolNone {
		return ""
	}
	if p.Protocol == ProtocolDirect {
		return "DIRECT"
	}
	if p.Username == "" && p.Password == "" {
		return string(p.Protocol) + "://" + p.Host + ":" + p.PortString()
	}
	if p.Username != "" && p.Password == "" {
		return string(p.Protocol) + "://" + p.Username + "@" + p.Host + ":" + p.PortString()
	}
	return string(p.Protocol) + "://" + p.Username + ":" + p.Password + "@" + p.Host + ":" + p.PortString()
}

// PortString returns the string representation of the port
func (p *Proxy[C]) PortString() string {
	return strconv.Itoa(int(p.Port))
}

// IsEmpty returns true if the proxy is empty
func (p *Proxy[C]) IsEmpty() bool {
	return p == nil || p.Protocol == ProtocolNone
}

// IsDirect returns true if the proxy is direct
func (p *Proxy[C]) IsDirect() bool {
	return p.Protocol == ProtocolDirect
}

// HasHttpClient returns true if the proxy has http client
func (p *Proxy[C]) HasHttpClient() bool {
	if p.httpClientCreator == nil && (p.httpClient == nil || p.httpClient.Len() == 0) {
		return false
	}
	return true
	// return !isNil[C](p.httpClient)
}

// GetHttpClient returns the http client for the specified key
//
// Will create a http client with the creator function if it does not exist
func (p *Proxy[C]) GetHttpClient(key ...string) C {
	var def C
	if len(key) == 0 {
		key = []string{""}
	}
	if p.httpClientCreator == nil && (p.httpClient == nil || p.httpClient.Len() == 0) {
		return def
	}
	if p.httpClient == nil {
		p.httpClient = shardmap.New[string, C](0)
	}
	hc, ok := p.httpClient.Get(key[0])
	if !ok {
		var err error
		hc, err = p.httpClientCreator(p)
		if err != nil {
			err = errors.Wrap(err, "failed to create http client")
			panic(err)
		}
		p.SetHttpClient(hc, key...)
	}
	if p.Rotating {
		switch v := any(hc).(type) {
		case tls_client.HttpClient:
			_ = v.SetProxy(p.String())
		case *http.Client:
			u, err := url.Parse(p.String())
			if err == nil {
				v.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
			}
		}
	}
	return hc
}

// SetHttpClient sets the http client
func (p *Proxy[C]) SetHttpClient(client C, key ...string) *Proxy[C] {
	if len(key) == 0 {
		key = []string{""}
	}
	if p.httpClient == nil {
		p.httpClient = shardmap.New[string, C](0)
	}
	p.httpClient.Set(key[0], client)
	return p
}
