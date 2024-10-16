package ve

import (
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/proxstore"
	tls_client "github.com/bogdanfinn/tls-client"
)

type VE struct {
	ps                *proxstore.ProxStore[tls_client.HttpClient]
	proxy             *proxstore.Proxy[tls_client.HttpClient]
	cj                *CookieJar
	config            *config.Config
	shapeSolverClient tls_client.HttpClient
}

func New(
	cfg *config.Config, proxyStore *proxstore.ProxStore[tls_client.HttpClient],
	proxy *proxstore.Proxy[tls_client.HttpClient],
) *VE {
	if proxy == nil {
		proxy = proxyStore.Next()
	}
	cj, _ := NewCookieJar(&CookieJarOptions{Options: nil})
	shapeSolverClient, _ := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithTimeoutSeconds(60),
	)
	return &VE{
		ps:                proxyStore,
		proxy:             proxy,
		cj:                cj,
		config:            cfg,
		shapeSolverClient: shapeSolverClient,
	}
}
