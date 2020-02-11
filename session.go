package cpanel

import (
	"fmt"
	"net/url"
)

type WhmAPI struct {
	hostname string
	token    string
}

// TODO figure out how to standardize ssh.Session and exec.Command

func (a WhmAPI) GenerateURL(endpoint string) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     a.hostname + ":2087",
		Path:     fmt.Sprintf("/%s/json-api/%s", a.token, endpoint),
		RawQuery: "?api.version=1",
	}
}
