package cpanel

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
)

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
	url, err := a.GenerateURL("listaccts")
	if err != nil {
		return []string{}, err
	}
	params := url.Query()
	params.Add("want", "user")
	url.RawQuery = params.Encode()

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return []string{}, err
	}
	response, err := a.client.Do(request)
	if err != nil {
		return []string{}, err
	}

	var outputData struct {
		Data struct {
			AccountList []struct {
				Username string `json:"user"`
			} `json:"acct"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&outputData)
	if err != nil {
		return []string{}, err
	}

	userlist := []string{}
	for _, account := range outputData.Data.AccountList {
		userlist = append(userlist, account.Username)
	}
	return userlist, nil
}

func (a *WhmAPI) ListResellers() ([]string, error) {
	var outputData struct {
		Data struct {
			Resellers []string `json:"reseller"`
		} `json:"data"`
	}
	queryParams := url.Values{}
	queryParams.Add("want", "user")

	err := a.Call(
		http.MethodGet,
		"listresellers",
		queryParams,
		&outputData,
	)

	if err != nil {
		return []string{}, err
	}

	return outputData.Data.Resellers, nil
}
