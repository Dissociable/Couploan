package proxstore

import "errors"

var (
	ErrInvalidProxyLine = errors.New("invalid proxy line")
	ErrInvalidProtocol  = errors.New("invalid protocol")
)
