package util

import "github.com/pkg/errors"

var (
	ErrNilProxy           = errors.New("proxy is nil, there must be at least a proxy with Direct protocol in case you don't have proxies")
	ErrInvalidCheckResult = errors.New("invalid check result")
)
