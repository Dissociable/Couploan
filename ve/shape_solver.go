package ve

import (
	"context"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
)

func (ve *VE) SolveShape(ctx context.Context) (err error) {
	if ve.proxy == nil || ve.proxy.IsDirect() {
		err = errors.New("no proxy is set")
		return
	}

	rb := url.Values{}
	rb.Set("url", "https://ve.cbi.ir/DefaultVE.aspx")
	rb.Set("proxy", ve.proxy.String())

	_, body, err := NewRequest(ve, "POST", ve.config.ShapeSolver.URL+"/api/v1/f5").
		SetContext(ctx).
		SetClient(ve.shapeSolverClient).
		SetProxy(ve.ps.Direct()).
		SetRetry().
		SetHeaders(
			map[string][]string{
				"Content-Type":    {"application/x-www-form-urlencoded"},
				"Authorization":   {"Bearer " + ve.config.ShapeSolver.ApiKey},
				"User-Agent":      {"Couploan"},
				"Accept-Encoding": {"gzip, deflate, br"},
			},
		).
		SetBody(strings.NewReader(rb.Encode())).
		SetMaxRetries(3).
		Do()
	if err != nil {
		err = errors.Wrap(err, "failed to solve shape")
		return
	}

	if !gjson.Valid(body) {
		err = fmt.Errorf("invalid json returned by ShapeSolver: %s", body)
		return
	}

	j := gjson.Parse(body)
	if !j.Get("success").Bool() {
		err = fmt.Errorf("failed to solve shape, un-successful response: %s", j.Get("message").String())
		return
	}

	if !j.Get("passed").Bool() {
		err = fmt.Errorf("failed to pass shape: %s", body)
		return
	}

	ck := j.Get("cookies").Array()

	for _, c := range ck {
		name := c.Get("name").String()
		value := c.Get("value").String()

		ve.cj.SetCookies(
			&url.URL{Scheme: "https", Host: "ve.cbi.ir", Path: "/"},
			[]*http.Cookie{{Name: name, Value: value}},
		)
	}
	return
}
