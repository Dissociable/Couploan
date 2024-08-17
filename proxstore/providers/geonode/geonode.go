package geonode

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"net/http"
)

type Service string

const (
	ServicePremiumResidential Service = "RESIDENTIAL-PREMIUM"
	ServicePremiumUnmetered   Service = "RESIDENTIAL-UNMETERED"
	ServiceSharedDatacenter   Service = "SHARED-DATACENTER"
)

type BasicParams struct {
	ServiceType Service
	Username    string
	Password    string
}

type ReleasePayload struct {
	Data       []ReleasePayloadData `json:"data,omitempty"`
	ReleaseAll bool                 `json:"releaseAll"`
}

type ReleasePayloadData struct {
	Port      int    `json:"port"`
	SessionId string `json:"sessionId,omitempty"`
}

func Release(basic BasicParams, all bool, data ...ReleasePayloadData) (bool, error) {
	p := ReleasePayload{Data: data, ReleaseAll: all}
	pm, err := json.Marshal(p)
	if err != nil {
		err = errors.Wrap(err, "failed to marshal payload")
		return false, err
	}
	r, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("https://monitor.geonode.com/sessions/release/%s", basic.ServiceType),
		bytes.NewReader(pm),
	)
	if err != nil {
		err = errors.Wrap(err, "failed to create request")
		return false, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.SetBasicAuth(basic.Username, basic.Password)

	resp, err := hc.Do(r)
	if err != nil {
		err = errors.Wrap(err, "failed to do request")
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}
