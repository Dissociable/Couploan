package geonode

import (
	"net/http"
	"time"
)

var (
	hc = http.Client{Timeout: 10 * time.Second}
)
