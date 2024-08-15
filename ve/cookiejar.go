package ve

import (
	cookiejar "github.com/Dissociable/persistent-cookiejar"
	http "github.com/bogdanfinn/fhttp"
	"github.com/pkg/errors"
	http2 "net/http"
	"net/url"
)

type CookieJarOptions struct {
	*cookiejar.Options
}

type CookieJar struct {
	Jar *cookiejar.Jar
}

func NewCookieJar(options *CookieJarOptions) (*CookieJar, error) {
	r, err := cookiejar.New(options.Options)
	if err != nil {
		err = errors.Wrap(err, "failed to create cookie Jar")
		return nil, err
	}
	return &CookieJar{
		Jar: r,
	}, nil
}

func (c CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	var transformed []*http2.Cookie

	for _, cookie := range cookies {
		transformed = append(
			transformed, &http2.Cookie{
				Name:       cookie.Name,
				Value:      cookie.Value,
				Path:       cookie.Path,
				Domain:     cookie.Domain,
				Expires:    cookie.Expires,
				RawExpires: cookie.RawExpires,
				MaxAge:     cookie.MaxAge,
				Secure:     cookie.Secure,
				HttpOnly:   cookie.HttpOnly,
				SameSite:   http2.SameSite(cookie.SameSite),
				Raw:        cookie.Raw,
				Unparsed:   cookie.Unparsed,
			},
		)
	}

	c.Jar.SetCookies(u, transformed)
}

func (c CookieJar) Cookies(u *url.URL) []*http.Cookie {
	r := c.Jar.Cookies(u)

	var transformed []*http.Cookie
	for _, cookie := range r {
		transformed = append(
			transformed, &http.Cookie{
				Name:       cookie.Name,
				Value:      cookie.Value,
				Path:       cookie.Path,
				Domain:     cookie.Domain,
				Expires:    cookie.Expires,
				RawExpires: cookie.RawExpires,
				MaxAge:     cookie.MaxAge,
				Secure:     cookie.Secure,
				HttpOnly:   cookie.HttpOnly,
				SameSite:   http.SameSite(cookie.SameSite),
				Raw:        cookie.Raw,
				Unparsed:   cookie.Unparsed,
			},
		)
	}
	return transformed
}

// Ensure interface compatibility
var _ http.CookieJar = (*CookieJar)(nil)
