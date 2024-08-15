package util

import (
	"context"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/pkg/errors"
	"io"
	"maps"
)

var UserAgent = UserAgentChrome
var UserAgentChrome = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
var GoogleBotUserAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

var DefaultGetHeaders = DefaultChromeHeaders

var DefaultChromeHeaders = http.Header{
	"sec-ch-ua":                 {"\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\""},
	"sec-ch-ua-mobile":          {"?0"},
	"sec-ch-ua-platform":        {"\"Windows\""},
	"upgrade-insecure-requests": {"1"},
	"user-agent":                {UserAgent},
	"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"sec-fetch-site":            {"cross-site"},
	"sec-fetch-mode":            {"navigate"},
	"sec-fetch-user":            {"?1"},
	"sec-fetch-dest":            {"document"},
	"referer":                   {"https://www.google.com/"},
	"accept-language":           {"en-US,en;q=0.9"},
	"priority":                  {"u=0, i"},
	"accept-encoding":           {"gzip, deflate, br"},
	http.HeaderOrderKey: {
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"referer",
		"accept-language",
		"priority",
		"accept-encoding",
	},
}

var DefaultGoogleBotHeaders = http.Header{
	"Accept-Language": {"en-US"},
	"Cache-Control":   {"no-cache"},
	"Connection":      {"keep-alive"},
	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"From":            {"googlebot(at)googlebot.com"},
	"User-Agent":      {GoogleBotUserAgent},
	"Accept-Encoding": {"gzip,deflate,br"},
	http.HeaderOrderKey: {
		"Accept-Language",
		"Cache-Control",
		"Connection",
		"Accept",
		"From",
		"User-Agent",
		"Accept-Encoding",
	},
}

func GetRequest(ctx context.Context, client tls_client.HttpClient, r *http.Request) (
	resp *http.Response, body string, err error,
) {
	r = r.WithContext(ctx)
	resp, err = client.Do(r)
	if err != nil {
		err = errors.Wrap(err, "failed to get request")
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read response body")
		return
	}

	body = string(readBytes)
	return
}

// BuildGetRequest creates a GET request
func BuildGetRequest(url string, headers http.Header) (r *http.Request, err error) {
	if headers == nil {
		headers = DefaultGetHeaders
	}
	headers = maps.Clone(headers)
	r, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to create http request")
		return
	}

	r.Header = headers
	return
}

// BuildPostRequest creates a POST request
func BuildPostRequest(url string, headers http.Header, requestBody io.Reader) (r *http.Request, err error) {
	if headers == nil {
		headers = DefaultGetHeaders
	}
	headers = maps.Clone(headers)
	r, err = http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		err = errors.Wrap(err, "failed to create http request")
		return
	}

	r.Header = headers
	return
}
