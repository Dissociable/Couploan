package ve

import (
	"context"
	"strings"
)

func (ve *VE) IP(ctx context.Context) (ip string, err error) {
	_, body, err := NewRequest(ve, "GET", "https://api.ipify.org").
		SetContext(ctx).
		SetGetClientFunc(getClientFunc).
		SetProxy(ve.proxy).SetRetry().
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
