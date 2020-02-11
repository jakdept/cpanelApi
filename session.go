package cpanel

import (
	"fmt"
	"net/url"
)

type WhmAPI struct {
	hostname   string
	sessionKey string
}

func (a WhmAPI) GenerateURL(endpoint string) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     a.hostname + ":2087",
		Path:     fmt.Sprintf("/%s/json-api/%s", a.sessionKey, endpoint),
		RawQuery: "?api.version=1",
	}
}
