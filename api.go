package cpanel

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

type WhmAPI struct {
	hostname *string
	token    *string
	client   *http.Client
}

// TODO figure out how to standardize ssh.Session and exec.Command

func (a *WhmAPI) GenerateURL(endpoint string) (url.URL, error) {
	if a.token == nil {
		return url.URL{}, errors.New("whm api endpoint not activated")
	}
	return url.URL{
		Scheme:   "https",
		Host:     *a.hostname + ":2087",
		Path:     fmt.Sprintf("%s/json-api/%s", *a.token, endpoint),
		RawQuery: "api.version=1",
	}, nil
}

func NewWhmApi(hostname string) (*WhmAPI, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return &WhmAPI{}, err
	}

	return &WhmAPI{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, //nolint
				},
			},
			Jar: jar,
		},
	}, nil
}
