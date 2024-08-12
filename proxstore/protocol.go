package proxstore

import "strings"

type Protocol string

var (
	ProtocolNone    = Protocol("none")
	ProtocolHttp    = Protocol("http")
	ProtocolHttps   = Protocol("https")
	ProtocolSocks4  = Protocol("socks4")
	ProtocolSocks4a = Protocol("socks4a")
	ProtocolSocks5  = Protocol("socks5")
	ProtocolSocks5h = Protocol("socks5h")
	ProtocolDirect  = Protocol("direct")
)

// ProtocolFromString converts string to Protocol
func ProtocolFromString(protocol string) Protocol {
	protocol = strings.ToLower(protocol)
	switch {
	case strings.HasPrefix(protocol, "https"):
		return ProtocolHttps
	case strings.HasPrefix(protocol, "http"):
		return ProtocolHttp
	case strings.HasPrefix(protocol, "socks4a"):
		return ProtocolSocks4a
	case strings.HasPrefix(protocol, "socks4"):
		return ProtocolSocks4
	case strings.HasPrefix(protocol, "socks5h"):
		return ProtocolSocks5h
	case strings.HasPrefix(protocol, "socks5"):
		return ProtocolSocks5
	case strings.HasPrefix(protocol, "direct"):
		return ProtocolDirect
	case strings.HasPrefix(protocol, "none"):
	default:
		return ProtocolNone
	}
	return ProtocolNone
}
