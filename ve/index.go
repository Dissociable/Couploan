package ve

import (
	"context"
	"strings"
)

func (ve *VE) Index(ctx context.Context) (ip string, err error) {
	// ve.cj.SetCookies(
	// 	&url.URL{Scheme: "https", Host: "ve.cbi.ir", Path: "/"},
	// 	[]*http.Cookie{
	// 		{
	// 			Name: "TSPD", Value: "1", Path: "/", Domain: "ve.cbi.ir", Secure: true,
	// 			Expires: time.Now().AddDate(0, 0, 1),
	// 		},
	// 	},
	// )
	_, body, err := NewRequest(ve, "GET", "https://ve.cbi.ir/DefaultVE.aspx").
		SetContext(ctx).
		SetGetClientFunc(getClientFunc).
		SetGetCookieJarFunc(getCookieJarFunc).
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
