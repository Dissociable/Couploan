package ve

import (
	"context"
	"github.com/Dissociable/Couploan/proxstore"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/pkg/errors"
	"strings"
)

func (ve *VE) IP(ctx context.Context) (ip string, err error) {
	_, body, err := NewRequest(ve, "GET", "https://api.ipify.org").
		SetContext(ctx).
		SetGetClientFunc(getClientFunc).
		SetProxy(ve.proxy).
		SetRetry().
		SetMaxRetries(3).
		Do()
	if err != nil {
		return "", err
	}

	body = strings.TrimSpace(body)
	if body != "" {
		ip = body
	}
	return
}

// CheckProxy checks if the proxy is valid and working,
// if proxy been working and good http response returned, it returns nil
func CheckProxy(ctx context.Context, proxy *proxstore.Proxy[tls_client.HttpClient]) (err error) {
	if proxy == nil {
		err = errors.New("no proxy is set")
		return
	}

	c, err := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithProxyUrl(proxy.String()),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithClientProfile(profiles.Chrome_124),
	)
	if err != nil {
		err = errors.Wrap(err, "failed to create http client")
		return
	}
	resp, body, err := NewRequest[any](nil, "GET", "https://ve.cbi.ir/DefaultVE.aspx").
		SetContext(ctx).
		SetClient(c).
		SetProxy(proxy).
		Do()

	if err != nil {
		err = errors.Wrap(err, "failed to check proxy")
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		if strings.Contains(body, `window["bobcmn"]`) || strings.Contains(body, "ازدواج") {
			return nil
		} else {
			err = errors.New("proxy is not working, couldn't verify response body")
			return err
		}
	}

	return errors.New("proxy is not working, bad response status code returned")
}
