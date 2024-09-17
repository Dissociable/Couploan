package ve

import (
	"github.com/Dissociable/Couploan/proxstore"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/pkg/errors"
)

var (
	getClientFunc func(requester *Requester[*VE]) (tls_client.HttpClient, error)

	afterResponse func(
		requester *Requester[*VE], resp *http.Response, respBody *string, err error,
	) error

	retryCheck func(requester *Requester[*VE], resp *http.Response, respBody *string, err error) bool

	getCookieJarFunc func(requester *Requester[*VE]) (*CookieJar, error)
)

func init() {
	// Get HTTP SetClient
	getClientFunc = func(requester *Requester[*VE]) (tls_client.HttpClient, error) {
		if requester.GetProxy().HasHttpClient() {
			client := requester.GetProxy().GetHttpClient("ve")
			return client, nil
		} else {
			err := errors.New("no http client is set to the proxy")
			return nil, err
		}
	}

	afterResponse = nil

	retryCheck = func(
		requester *Requester[*VE], resp *http.Response, respBody *string, err error,
	) bool {
		if err != nil || resp == nil || resp.StatusCode != 200 {
			return true
		}
		return false
	}

	// Cookie
	getCookieJarFunc = func(requester *Requester[*VE]) (*CookieJar, error) {
		if requester.GetCookieJar() != nil {
			return requester.GetCookieJar(), nil
		} else {
			return requester.GetBase().cj, nil
		}
	}
}

// getProxy returns the last used proxy
func (ve *VE) getProxy() *proxstore.Proxy[tls_client.HttpClient] {
	return ve.ps.Next()
}
