package proxstore

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProxStoreLoadLine(t *testing.T) {
	p := New()
	noErrorProxies := []string{
		"http://username@127.0.0.1:8080",
		"http://username:password@127.0.0.1:8080",
		"http://127.0.0.1:8080",
		"socks5://127.0.0.1:8080",
	}
	for _, noErrorProxy := range noErrorProxies {
		err := p.LoadLine(noErrorProxy)
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		if !assert.False(t, p.Last().HasHttpClient()) {
			assert.FailNow(t, "http client should be nil")
		}
	}

	errorProxies := []string{
		"127.0.0.1",
		"http://127.0.0.1",
		"socks52://127.0.0.1",
		"hello",
		"hello:world",
		"http://hello:world",
		"socks5://username:password@hello:world",
	}
	for _, errorProxy := range errorProxies {
		err := p.LoadLine(errorProxy)
		if !assert.Error(t, err) {
			assert.FailNow(t, err.Error())
		}
		if !assert.False(t, p.Last().HasHttpClient()) {
			assert.FailNow(t, "http client should be nil")
		}
	}

	assert.Equal(t, len(noErrorProxies), p.Count())
}

func TestProxStoreWithOptionsLoadLine(t *testing.T) {
	optionChc := &OptionsCreateHttpClient[tls_client.HttpClient]{
		Creator: func(proxy *Proxy[tls_client.HttpClient]) (hc tls_client.HttpClient, err error) {
			options := []tls_client.HttpClientOption{
				tls_client.WithTimeoutSeconds(30),
				tls_client.WithClientProfile(profiles.Chrome_120),
				tls_client.WithNotFollowRedirects(),
			}
			if !proxy.IsEmpty() && !proxy.IsDirect() {
				options = append(options, tls_client.WithProxyUrl(proxy.String()))
			}
			hc, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
			return
		},
	}

	p := NewWithOptions[tls_client.HttpClient](DefaultOptions, optionChc)
	noErrorProxies := []string{
		"http://username@127.0.0.1:8080",
		"http://username:password@127.0.0.1:8080",
		"http://127.0.0.1:8080",
		"socks5://127.0.0.1:8080",
	}
	for _, noErrorProxy := range noErrorProxies {
		err := p.LoadLine(noErrorProxy)
		if !assert.NoError(t, err) {
			assert.FailNow(t, err.Error())
		}
		if !assert.True(t, p.Last().HasHttpClient()) {
			assert.FailNow(t, "http client should not be nil")
		}
	}

	errorProxies := []string{
		"127.0.0.1",
		"http://127.0.0.1",
		"socks52://127.0.0.1",
		"hello",
		"hello:world",
		"http://hello:world",
		"socks5://username:password@hello:world",
	}
	for _, errorProxy := range errorProxies {
		err := p.LoadLine(errorProxy)
		if !assert.Error(t, err) {
			assert.FailNow(t, err.Error())
		}
		if !assert.True(t, p.Last().HasHttpClient()) {
			assert.FailNow(t, "http client should not be nil")
		}
	}

	assert.Equal(t, p.Last().String(), p.Last().GetHttpClient().GetProxy())

	assert.Equal(t, len(noErrorProxies), p.Count())
}
