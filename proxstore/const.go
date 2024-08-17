package proxstore

type ProviderName string

type Provider struct {
	Name        ProviderName
	ServiceType string
	Username    string
	Password    string
}

const (
	ProviderNameGeoNode ProviderName = "geonode"
	ProviderNameNone    ProviderName = ""
)

var (
	DefaultOptions = &Options{AllowDirect: false, Provider: &Provider{Name: ProviderNameNone}}
)
