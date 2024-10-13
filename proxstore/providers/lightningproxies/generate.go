package lightningproxies

import (
	"fmt"
	"github.com/Dissociable/Couploan/proxstore"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type GenerateSettings struct {
	Username   string
	Password   string
	Zone       Zone
	Protocol   Protocol
	Region     *string
	State      *string
	City       *string
	ISP        *string // ISP is the ASN of the ISP, ASXXXX
	Rotating   bool
	StickyTime int // StickyTime in minutes
}

func (s *GenerateSettings) String() (result string, err error) {
	host, hostOk := ZoneHosts[s.Zone]
	if !hostOk {
		err = errors.New("Hosts map has no host for the specified zone")
		return
	}
	port := 0
	if s.Protocol == ProtocolHTTP {
		port = HttpPort
	} else if s.Protocol == ProtocolSOCKS5 {
		port = Socks5Port
	}
	var u []string
	u = append(u, s.Username)
	u = append(u, "zone", s.Zone.String())
	if s.Region != nil && len(*s.Region) > 0 {
		u = append(u, "region", strings.ToLower(*s.Region))
	}
	if s.State != nil && len(*s.State) > 0 {
		u = append(u, "st", strings.ToLower(*s.State))
	}
	if s.City != nil && len(*s.City) > 0 {
		u = append(u, "city", strings.ToLower(*s.City))
	}
	if s.ISP != nil && len(*s.ISP) > 0 {
		u = append(u, "asn", strings.ToLower(*s.ISP))
	}
	if !s.Rotating {
		u = append(u, "session", proxstore.RandomString(StickySessionLength), "sessTime", strconv.Itoa(s.StickyTime))
	}
	result = string(s.Protocol) + "://" +
		strings.Join(u, "-") + ":" +
		s.Password + "@" +
		host + ":" + strconv.Itoa(port)
	return
}

func GenerateProxies(s *GenerateSettings, count int) (proxies []string) {
	if s.Rotating {
		count = 1
	}
	for i := 0; i < count; i++ {
		proxy, err := s.String()
		if err != nil {
			fmt.Println("failed to generate proxy: " + err.Error())
			continue
		}
		proxies = append(proxies, proxy)
	}
	return
}
