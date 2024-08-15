package ve

import (
	"context"
	"github.com/Dissociable/Couploan/proxstore"
	"github.com/Dissociable/Couploan/ve/util"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	_ "github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
	"io"
)

type RequesterRetryCheck[C any] func(requester *Requester[C], resp *http.Response, respBody *string, err error) bool

func NewRequesterDefaultRetryCheck[C any]() RequesterRetryCheck[C] {
	f := func(
		requester *Requester[C], resp *http.Response, respBody *string, err error,
	) bool {
		if err != nil || resp == nil || resp.StatusCode != 200 {
			return true
		}
		return false
	}
	return f
}

type Requester[C any] struct {
	base             C
	ctx              context.Context
	client           tls_client.HttpClient
	getClientFunc    func(requester *Requester[C]) (tls_client.HttpClient, error)
	method           string
	link             string
	proxy            *proxstore.Proxy[tls_client.HttpClient]
	headers          http.Header
	body             io.Reader
	cookieJar        *CookieJar
	getCookieJarFunc func(requester *Requester[C]) (*CookieJar, error)
	retry            bool
	maxRetries       int
	retries          int
	retryCheck       func(requester *Requester[C], resp *http.Response, respBody *string, err error) bool
	afterResponse    func(requester *Requester[C], resp *http.Response, respBody *string, err error) error
}

func NewRequest[C any](base C, method string, link string) *Requester[C] {
	return &Requester[C]{
		base:       base,
		method:     method,
		link:       link,
		retry:      false,
		maxRetries: 3,
		retryCheck: NewRequesterDefaultRetryCheck[C](),
		ctx:        context.Background(),
	}
}

// SetClient sets the http client
func (r *Requester[C]) SetClient(client tls_client.HttpClient) *Requester[C] {
	r.client = client
	return r
}

// SetGetClientFunc sets the method that returns the http client
func (r *Requester[C]) SetGetClientFunc(
	getClientFunc func(requester *Requester[C]) (
		tls_client.HttpClient, error,
	),
) *Requester[C] {
	r.getClientFunc = getClientFunc
	return r
}

func (r *Requester[C]) SetRetry() *Requester[C] {
	r.retry = true
	return r
}

func (r *Requester[C]) SetMaxRetries(maxRetries int) *Requester[C] {
	r.maxRetries = maxRetries
	return r
}

// SetRetryCheck sets the method that determines when to retry
func (r *Requester[C]) SetRetryCheck(
	check func(
		requester *Requester[C], resp *http.Response, respBody *string, err error,
	) bool,
) *Requester[C] {
	r.retryCheck = check
	return r
}

func (r *Requester[C]) SetMethod(method string) *Requester[C] {
	r.method = method
	return r
}

func (r *Requester[C]) SetLink(link string) *Requester[C] {
	r.link = link
	return r
}

func (r *Requester[C]) SetProxy(proxy *proxstore.Proxy[tls_client.HttpClient]) *Requester[C] {
	r.proxy = proxy
	return r
}

func (r *Requester[C]) SetHeaders(headers http.Header) *Requester[C] {
	r.headers = headers
	return r
}

// SetBody sets the body of the request
//
// NOTE: DO NOT USE io.NopCloser(...)
func (r *Requester[C]) SetBody(body io.Reader) *Requester[C] {
	r.body = io.NopCloser(body)
	return r
}

func (r *Requester[C]) SetCookieJar(cookieJar *CookieJar) *Requester[C] {
	r.cookieJar = cookieJar
	return r
}

// SetGetCookieJarFunc sets the method that will be called to get the cookie jar
func (r *Requester[C]) SetGetCookieJarFunc(
	getCookieJarFunc func(requester *Requester[C]) (
		*CookieJar, error,
	),
) *Requester[C] {
	r.getCookieJarFunc = getCookieJarFunc
	return r
}

// SetAfterResponse sets the method that will be called after the response
//
// In case it returns error, no retries will be made, and the error will be returned by Do()
func (r *Requester[C]) SetAfterResponse(
	afterResponse func(
		requester *Requester[C], resp *http.Response, respBody *string, err error,
	) error,
) *Requester[C] {
	r.afterResponse = afterResponse
	return r
}

func (r *Requester[C]) SetContext(ctx context.Context) *Requester[C] {
	r.ctx = ctx
	return r
}

// Get Methods

func (r *Requester[C]) GetBase() C {
	return r.base
}

func (r *Requester[C]) GetClient() tls_client.HttpClient {
	return r.client
}

func (r *Requester[C]) GetClientFunc() func(requester *Requester[C]) (tls_client.HttpClient, error) {
	return r.getClientFunc
}

func (r *Requester[C]) GetProxy() *proxstore.Proxy[tls_client.HttpClient] {
	return r.proxy
}

func (r *Requester[C]) GetCookieJar() *CookieJar {
	return r.cookieJar
}

func (r *Requester[C]) GetLink() string {
	return r.link
}

func (r *Requester[C]) GetHeaders() http.Header {
	return r.headers
}

func (r *Requester[C]) GetBody() io.Reader {
	return r.body
}

func (r *Requester[C]) GetContext() context.Context {
	return r.ctx
}

func (r *Requester[C]) Do() (
	resp *http.Response, respBody string, err error,
) {
	if r.getCookieJarFunc != nil && r.cookieJar == nil {
		jar, err := r.getCookieJarFunc(r)
		if err != nil {
			err = errors.Wrap(err, "failed to get cookie jar")
			return resp, respBody, err
		}
		r.cookieJar = jar
	}
	if r.proxy == nil {
		err = util.ErrNilProxy
		return
	}
	if r.headers == nil {
		r.headers = util.DefaultGetHeaders
	}
	var req *http.Request
	if r.method == http.MethodGet {
		req, err = util.BuildGetRequest(r.link, r.headers)
		if err != nil {
			return
		}
	} else if r.method == http.MethodPost {
		req, err = util.BuildPostRequest(r.link, r.headers, r.body)
		if err != nil {
			return
		}
	}
	var client tls_client.HttpClient
	// If r.client is not nil and if proxy is not nil and proxy is not rotating, the re-use the client
	// otherwise, get the client again
	if r.client != nil && ((r.proxy != nil && (r.proxy.Rotating)) || r.proxy == nil) {
		client = r.client
	} else if r.getClientFunc != nil {
		client, err = r.getClientFunc(r)
		if err != nil {
			err = errors.Wrap(err, "failed to get http client")
			return
		}
		r.client = client
	} else if r.client != nil {
		client = r.client
	}
	if client == nil {
		err = errors.New("no http client is set")
		return
	}

	if r.cookieJar != nil {
		client.SetCookieJar(r.cookieJar)
	}
	resp, respBody, err = util.GetRequest(r.ctx, client, req)
	if r.afterResponse != nil {
		errAfterResponse := r.afterResponse(r, resp, &respBody, err)
		if errAfterResponse != nil {
			err = errors.Wrap(errAfterResponse, "afterResponse function errored")
			return
		}
	}

	if r.retry &&
		r.retryCheck != nil &&
		r.retryCheck(r, resp, &respBody, err) &&
		r.maxRetries > 0 &&
		r.maxRetries > r.retries {

		r.retries++
		return r.Do()
	}
	if err != nil {
		err = errors.Wrap(err, "failed to get request")
		return
	}
	return
}

// ReSetProxy re-sets the proxy to the requester and gets the client again
func (r *Requester[C]) ReSetProxy(proxy *proxstore.Proxy[tls_client.HttpClient]) (err error) {
	r.SetProxy(proxy)
	if r.GetClientFunc() != nil {
		newClient, err := r.GetClientFunc()(r)
		if err != nil {
			err = errors.Wrap(err, "unable to get new client after proxy ban")
			return err
		}
		r.SetClient(newClient)
		return nil
	}
	return nil
}
