package proxstore

import (
	"github.com/Dissociable/Couploan/proxstore/providers/geonode"
	"github.com/phuslu/shardmap"
	"github.com/pkg/errors"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"
)

type Options struct {
	AllowDirect bool // AllowDirect whether to allow direct connections for when there is no proxies loaded
	Provider    *Provider
}

type OptionsCreateHttpClient[C any] struct {
	// Creator is the function that creates the http client
	//
	// NOTE: If you want to change it after the creation, use [ProxStore.SetHttpClientCreator]
	Creator CreateHttpClientCreator[C]
}

type CreateHttpClientCreator[C any] func(proxy *Proxy[C]) (hc C, err error)

type ProxStore[C any] struct {
	proxies                *shardmap.Map[string, *Proxy[C]]
	options                *Options
	optionCreateHttpClient *OptionsCreateHttpClient[C]
	directProxy            *Proxy[C]
	rand                   *rand.Rand
	// index atomic number as counter for proxy index
	index *atomic.Int32
}

func New() *ProxStore[any] {
	direct := NewProxy[any]("", 0, ProtocolDirect)
	index := &atomic.Int32{}
	index.Store(-1)
	return &ProxStore[any]{
		proxies:                shardmap.New[string, *Proxy[any]](0),
		options:                DefaultOptions,
		optionCreateHttpClient: nil,
		directProxy:            direct,
		rand:                   rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		index:                  index,
	}
}

func NewWithOptions[C any](options *Options, optionCreateHttpClient *OptionsCreateHttpClient[C]) *ProxStore[C] {
	if options == nil {
		options = DefaultOptions
	}
	var direct *Proxy[C]
	if options.AllowDirect {
		direct = NewProxy[C]("", 0, ProtocolDirect)
		if optionCreateHttpClient != nil && optionCreateHttpClient.Creator != nil {
			direct.SetHttpClientCreator(optionCreateHttpClient.Creator)
		}
	}
	index := &atomic.Int32{}
	index.Store(-1)
	return &ProxStore[C]{
		proxies:                shardmap.New[string, *Proxy[C]](0),
		options:                options,
		optionCreateHttpClient: optionCreateHttpClient,
		directProxy:            direct,
		rand:                   rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		index:                  index,
	}
}

// SetHttpClientCreator allows you to set a custom http client creator
func (p *ProxStore[C]) SetHttpClientCreator(creator CreateHttpClientCreator[C]) {
	p.optionCreateHttpClient = &OptionsCreateHttpClient[C]{
		Creator: creator,
	}
}

func (p *ProxStore[C]) LoadFromFile() (err error) {
	return
}

func (p *ProxStore[C]) LoadProxy(proxy *Proxy[C]) (err error) {
	if proxy.IsEmpty() {
		return ErrInvalidProxyLine
	}
	p.loadProxy(proxy)
	return
}

func (p *ProxStore[C]) loadProxy(proxy *Proxy[C]) {
	if p.optionCreateHttpClient != nil && p.optionCreateHttpClient.Creator != nil && !proxy.HasHttpClient() {
		if proxy.httpClientCreator == nil {
			proxy.httpClientCreator = p.optionCreateHttpClient.Creator
		}
	}
	p.proxies.Set(proxy.String(), proxy)
}

func (p *ProxStore[C]) LoadLine(line string, protocol ...Protocol) (err error) {
	line = strings.TrimSpace(line)
	lineSplit := strings.Split(line, ":")
	if len(lineSplit) < 2 {
		return ErrInvalidProxyLine
	}
	hasProtocol := false
	// check whether the line has protocol
	if len(line) > 10 {
		loweredFirstElement := strings.ToLower(line[0:10])
		if strings.HasPrefix(loweredFirstElement, "http://") ||
			strings.HasPrefix(loweredFirstElement, "https://") ||
			strings.HasPrefix(loweredFirstElement, "socks4://") ||
			strings.HasPrefix(loweredFirstElement, "socks4a://") ||
			strings.HasPrefix(loweredFirstElement, "socks5://") ||
			strings.HasPrefix(loweredFirstElement, "socks5h://") {
			hasProtocol = true
		}
	}
	// line has no credentials
	if hasProtocol {
		// Line has protocol
		err = p.loadLineWithProtocol(line, lineSplit)
		if err != nil {
			err = errors.Wrap(err, "failed to load line with protocol")
			return
		}
	} else {
		// If protocol is None, then return error as line has no protocol
		if len(protocol) == 0 {
			return ErrInvalidProxyLine
		}
		// Line has no protocol
		err = p.loadLineWithoutProtocol(line, lineSplit, protocol[0])
		if err != nil {
			err = errors.Wrap(err, "failed to load line without protocol")
			return
		}
	}
	return
}

// loadLineWithProtocol loads a line with protocol, e.g., http://127.0.0.1:8080 or socks5://127.0.0.1:8080
// or with credential: http://username:password@127.0.0.1:8080 or without password: http://username@127.0.0.1:8080
//
// lineSplit must be the line split via ":"
func (p *ProxStore[C]) loadLineWithProtocol(line string, lineSplit []string) (err error) {
	proxy, err := ParseLineWithProtocol[C](line, lineSplit)
	if err != nil {
		return err
	}
	if proxy.IsEmpty() {
		return ErrInvalidProxyLine
	}
	p.loadProxy(proxy)
	return
}

// ParseLineWithProtocol parses a line with protocol, e.g., http://127.0.0.1:8080 or socks5://127.0.0.1:8080
// or with credential: http://username:password@127.0.0.1:8080 or without password: http://username@127.0.0.1:8080
//
// lineSplit must be the line split via ":"
func ParseLineWithProtocol[C any](line string, lineSplit []string) (proxy *Proxy[C], err error) {
	protocol := lineSplit[0]
	if len(lineSplit) < 3 {
		err = ErrInvalidProxyLine
		return
	}
	var host, username, password string
	var port uint16
	// Line has credentials
	// parse username
	username = lineSplit[1]
	// remove the double / from the starting of username since the lineSplit been done via splitting by ":" delimiter
	username = strings.TrimPrefix(username, "//")
	// line has credential without a password
	if len(lineSplit) < 4 {
		// split username and host from username variable
		if strings.Contains(username, "@") {
			usernameSplit := strings.Split(username, "@")
			usernameSplitFirstElement := ""
			if len(usernameSplit) > 2 {
				usernameSplitFirstElement = strings.Join(usernameSplit[0:len(usernameSplit)-2], "@")
			} else {
				usernameSplitFirstElement = strings.Join(usernameSplit[0:len(usernameSplit)-1], "@")
			}
			username = usernameSplitFirstElement
			host = usernameSplit[len(usernameSplit)-1]
		} else {
			host = username
			username = ""
		}
	}
	// line has credential with a password too
	if len(lineSplit) >= 4 {
		// parse password
		password = lineSplit[2]
		// split password and host from password variable
		if strings.Contains(password, "@") {
			passwordSplit := strings.Split(password, "@")
			passwordSplitFirstElement := ""
			if len(passwordSplit) > 2 {
				passwordSplitFirstElement = strings.Join(passwordSplit[0:len(passwordSplit)-2], "@")
			} else {
				passwordSplitFirstElement = strings.Join(passwordSplit[0:len(passwordSplit)-1], "@")
			}
			password = passwordSplitFirstElement
			host = passwordSplit[len(passwordSplit)-1]
		}
	}
	// parse port
	portStr := lineSplit[len(lineSplit)-1]
	// portStr to uint16
	portInt64, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		err = errors.Wrap(err, "failed to parse port")
		return
	}
	port = uint16(portInt64)
	protocolCasted := ProtocolFromString(protocol)
	if protocolCasted == ProtocolNone {
		err = ErrInvalidProtocol
		return
	}
	proxy = NewProxyWithCredential[C](host, port, protocolCasted, username, password)
	return
}

// loadLineWithoutProtocol loads a line without a protocol, e.g., 127.0.0.1:8080
// or with credential: 127.0.0.1:8080:username:password or without: 127.0.0.1:8080:username
//
// lineSplit must be the line split via ":"
func (p *ProxStore[C]) loadLineWithoutProtocol(line string, lineSplit []string, protocol Protocol) (err error) {
	proxy, err := ParseLineWithoutProtocol[C](line, lineSplit, protocol)
	if err != nil {
		return err
	}
	if proxy.IsEmpty() {
		return ErrInvalidProxyLine
	}
	p.loadProxy(proxy)
	return
}

// ParseLineWithoutProtocol Parses a line without a protocol, e.g., 127.0.0.1:8080
// or with credential: 127.0.0.1:8080:username:password or without: 127.0.0.1:8080:username
//
// lineSplit must be the line split via ":"
func ParseLineWithoutProtocol[C any](line string, lineSplit []string, protocol Protocol) (proxy *Proxy[C], err error) {
	if len(lineSplit) < 2 {
		err = ErrInvalidProxyLine
		return
	}
	var host, username, password string
	var port uint16
	// Line has credentials
	// parse host
	host = lineSplit[0]
	// parse port
	portStr := lineSplit[1]
	// portStr to uint16
	portInt64, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		err = errors.Wrap(err, "failed to parse port")
		return
	}
	port = uint16(portInt64)
	if len(lineSplit) > 2 {
		// parse username
		username = lineSplit[2]
		if len(lineSplit) > 3 {
			// parse password
			password = lineSplit[3]
		}
	}
	proxy = NewProxyWithCredential[C](host, port, protocol, username, password)
	return
}

// Count returns the number of proxies
func (p *ProxStore[C]) Count() int {
	return p.proxies.Len()
}

// Last returns the last proxy, nil if there are no proxies and Direct is not allowed, Direct otherwise.
func (p *ProxStore[C]) Last() *Proxy[C] {
	if p.proxies.Len() == 0 {
		return p.Direct()
	}
	var last *Proxy[C]
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			last = value
			return true
		},
	)
	return last
}

// First returns the first proxy, nil if there are no proxies and Direct is not allowed, Direct otherwise.
func (p *ProxStore[C]) First() *Proxy[C] {
	if p.proxies.Len() == 0 {
		return p.Direct()
	}
	var first *Proxy[C]
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			first = value
			return false
		},
	)
	return first
}

// Next returns the next proxy, nil if there are no proxies and Direct is not allowed, Direct otherwise.
func (p *ProxStore[C]) Next() *Proxy[C] {
	count := p.proxies.Len()
	if count == 0 {
		return p.Direct()
	}
	if count == 1 {
		return p.First()
	}
	index := p.index.Load()
	if int(index) >= p.proxies.Len() {
		index = 0
		p.index.Store(-1)
	}
	var prox *Proxy[C]
	counter := int32(-1)
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			counter++
			if counter == index+1 {
				prox = value
				return false
			}
			return true
		},
	)
	if prox == nil {
		prox = p.First()
		p.index.Store(-1)
	} else {
		p.index.Store(counter)
	}
	return prox
}

// ProxyAt returns the proxy at the given index, nil if there are no proxies and Direct is not allowed, Direct otherwise.
func (p *ProxStore[C]) ProxyAt(index int) *Proxy[C] {
	count := p.proxies.Len()
	if count == 0 {
		return p.Direct()
	}
	if index < 0 || index >= p.proxies.Len() {
		index = 0
	}
	var prox *Proxy[C]
	counter := -1
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			counter++
			if counter == index+1 {
				prox = value
				return false
			}
			return true
		},
	)
	if prox == nil {
		prox = p.First()
	}
	return prox
}

// Random returns a random proxy, nil if there are no proxies and Direct is not allowed, Direct otherwise.
func (p *ProxStore[C]) Random() *Proxy[C] {
	count := p.proxies.Len()
	if count == 0 {
		return p.Direct()
	}
	var prox *Proxy[C]
	if count == 1 {
		prox = p.First()
		return prox
	}
	randNumber := p.rand.IntN(count)
	counter := -1
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			counter++
			if counter == randNumber {
				prox = value
				return false
			}
			return true
		},
	)
	if prox == nil {
		prox = p.First()
	}
	return prox
}

// Direct returns the direct proxy that's been initialized via ProxStore initialization.
//
// Returns nil if [Options.AllowDirect] is false
func (p *ProxStore[C]) Direct() *Proxy[C] {
	return p.directProxy
}

// IsLastIndex returns if the index is the last index
func (p *ProxStore[C]) IsLastIndex() bool {
	lastIndex := p.proxies.Len() - 1
	if lastIndex < 1 {
		return true
	}
	index := int(p.index.Load())
	return index >= lastIndex
}

// GetIndex returns the index
func (p *ProxStore[C]) GetIndex() int {
	return int(p.index.Load())
}

// GetProxyIndex returns the index of the proxy in the ProxStore
//
// returns -1 if not found
func (p *ProxStore[C]) GetProxyIndex(proxy *Proxy[C]) int {
	if proxy == nil {
		return -1
	}

	index := -1
	p.proxies.Range(
		func(key string, value *Proxy[C]) bool {
			index++
			return value != proxy
		},
	)
	return index
}

// ReleaseProxy releases the proxy, if proxy is nil, releases all
func (p *ProxStore[C]) ReleaseProxy(proxy *Proxy[C]) (bool, error) {
	if proxy != nil {
		proxy.reloadIp.Store(true)
	}
	if (proxy == nil || !proxy.HasProvider()) && (p.options == nil || p.options.Provider == nil) {
		return false, nil
	}
	var provider *Provider
	if proxy != nil && proxy.HasProvider() {
		provider = proxy.provider
	} else {
		provider = p.options.Provider
	}
	var data []geonode.ReleasePayloadData
	if proxy != nil {
		d := geonode.ReleasePayloadData{Port: int(proxy.Port)}
		// extract sessionId from username, sessionId is between session- and the next -
		if strings.Contains(proxy.Username, "session-") {
			sessionId := strings.Split(proxy.Username, "session-")[1]
			sessionId = strings.Split(sessionId, "-")[0]
			d.SessionId = sessionId
		}
		data = append(data, d)
	}
	switch provider.Name {
	case ProviderNameGeoNode:
		return geonode.Release(
			geonode.BasicParams{
				ServiceType: geonode.Service(provider.ServiceType),
				Username:    provider.Username,
				Password:    provider.Password,
			},
			proxy == nil,
			data...,
		)
	}
	return false, nil
}
