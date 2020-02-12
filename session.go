package cpanel

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type WhmAPI struct {
	hostname string
	token    string
	client   *http.Client
}

// TODO figure out how to standardize ssh.Session and exec.Command

func (a *WhmAPI) GenerateURL(endpoint string) (*http.Client, url.URL) {
	return a.client,
		url.URL{
			Scheme:   "https",
			Host:     a.hostname + ":2087",
			Path:     fmt.Sprintf("/%s/json-api/%s", a.token, endpoint),
			RawQuery: "api.version=1",
		}
}

// Account represents a cPanel account
type Account struct {
	Username     string
	Domain       string
	ContactEmail string

	Reseller      string
	HomePartition string
	Shell         string
	Package       string
	Theme         string

	BackupsEnabled        bool
	Suspended             bool
	Locked                bool
	OutgoingMailSuspended bool
	OutgoingMailHold      bool

	MailboxFormat            int8
	MaxDeferPrecent          string
	MinDeferBeforeProtection string
	MaxEmailPerHour          string

	MainIPv4 net.IP
	MainIPv6 net.IP

	EmailQuotaLimit string
	MaxAddons       string
	MaxFtp          string
	MaxMailingLists string
	MaxParked       string
	MaxPop          string
	MaxDatabases    string
	MaxSubdomains   string

	DiskLimit  int
	DiskUsed   int
	InodeLimit int
	InodeUsed  int
}

func (a *WhmAPI) ListAccounts() ([]string, error) {
	client, url := a.GenerateURL("listaccts")
	params := url.Query()
	params.Add("want", "user")
	url.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return []string{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return []string{}, err
	}

	var outputData struct {
		AccountList []struct {
			Username string `json:"user"`
		} `json:"acct"`
	}

	err = json.NewDecoder(response.Body).Decode(&outputData)
	if err != nil {
		return []string{}, err
	}

	userlist := []string{}
	for _, account := range outputData.AccountList {
		userlist = append(userlist, account.Username)
	}
	return userlist, nil
}

func (a *WhmAPI) ListResellers() ([]string, error) {
	client, url := a.GenerateURL("listresellers")
	params := url.Query()
	params.Add("want", "user")
	url.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return []string{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return []string{}, err
	}

	var outputData struct {
		Resellers []string `json:"reseller"`
	}

	err = json.NewDecoder(response.Body).Decode(&outputData)
	if err != nil {
		return []string{}, err
	}

	return outputData.Resellers, nil
}
