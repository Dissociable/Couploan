package proxstore

var (
	DefaultOptions = &Options{}
	ProxyDirect    = NewProxy[any]("", 0, ProtocolDirect)
)
