package lightningproxies

type Protocol string

const (
	ProtocolHTTP   Protocol = "http"
	ProtocolSOCKS5 Protocol = "socks5"

	// HttpPort is the http proxy port
	HttpPort = 9999
	// Socks5Port is the socks5 proxy port
	Socks5Port = 9999

	// StickySessionLength Sticky Session Length
	StickySessionLength = 12
)

type Zone string

const (
	ZoneResidential Zone = "resi"
)

var (
	ZoneHosts = map[Zone]string{
		ZoneResidential: "resi-as.lightningproxies.net",
	}
)

func (z Zone) String() string {
	return string(z)
}
