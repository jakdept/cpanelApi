package cpanel

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func NewWhmApi(hostname string) (*WhmAPI, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return &WhmAPI{}, err
	}

	// trim off any port numbers and change it to WHM's port
	hostname = net.JoinHostPort(strings.Split(hostname, ":")[0], "2087")

	return &WhmAPI{
		hostname: &hostname,
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

type WhmAPI struct {
	hostname *string
	token    *string
	client   *http.Client
}

func (a *WhmAPI) Call(
	method string,
	endpoint string,
	args url.Values,
	out interface{},
) error {

	args["api.version"] = []string{"1"}
	url := url.URL{
		Scheme:   "https",
		Host:     *a.hostname,
		RawQuery: args.Encode(),
	}
	if a.token != nil {
		url.Path = fmt.Sprintf("%s/json-api/%s", *a.token, endpoint)
	} else {
		url.Path = fmt.Sprintf("/json-api/%s", endpoint)
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	response, err := a.client.Do(request)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&out)
	if err != nil {
		return err
	}
	return nil
}

type CpBool bool

func (b *CpBool) String() string {
	if *b {
		return "1"
	}
	return "0"
}

func (b *CpBool) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *CpBool) UnmarshalJSON(v []byte) error {
	s := string(v)
	s = strings.Trim(s, "\"")
	if s == "1" {
		*b = true
	} else {
		*b = false
	}
	return nil
}

type IntLimit struct {
	value     int
	unlimited bool
}

func (l *IntLimit) String() string {
	if l.unlimited {
		return "unlimited"
	}
	return strconv.Itoa(l.value)
}

func (l *IntLimit) MarshalJSON() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *IntLimit) UnmarshalJSON(v []byte) error {
	s := string(v)
	s = strings.Trim(s, "\"")
	if s == "unlimited" || s == "0" {
		l.unlimited = true
		l.value = 0
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("value is not a limit: %w", err)
	}
	l.value = i
	l.unlimited = false
	return nil
}

type FloatLimit struct {
	value     float64
	unlimited bool
}

func (l *FloatLimit) String() string {
	if l.unlimited {
		return "unlimited"
	}
	return strconv.FormatFloat(l.value, 'f', 2, 64)
}

func (l *FloatLimit) MarshalJSON() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *FloatLimit) UnmarshalJSON(v []byte) error {
	s := string(v)
	s = strings.Trim(s, "\"")
	if s == "unlimited" || s == "0" {
		l.unlimited = true
		l.value = 0
		return nil
	}

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("value is not a limit: %w", err)
	}
	l.value = i
	l.unlimited = false
	return nil
}
