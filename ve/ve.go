package ve

import (
	"github.com/Dissociable/Couploan/proxstore"
	tls_client "github.com/bogdanfinn/tls-client"
)

type VE struct {
	ps    *proxstore.ProxStore[tls_client.HttpClient]
	proxy *proxstore.Proxy[tls_client.HttpClient]
	cj    *CookieJar
}

func New(proxyStore *proxstore.ProxStore[tls_client.HttpClient], proxy *proxstore.Proxy[tls_client.HttpClient]) *VE {
	if proxy == nil {
		proxy = proxyStore.Next()
	}
	cj, _ := NewCookieJar(&CookieJarOptions{Options: nil})
	return &VE{
		ps:    proxyStore,
		proxy: proxy,
		cj:    cj,
	}
}
