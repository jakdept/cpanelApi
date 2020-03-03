package cpanel

import (
	"net"
	"net/http"
	"net/url"
)

// Account represents a cPanel account
type Account struct {
	Username      string `json:"user,omitempty"`
	PrimaryDomain string `json:"domain,omitempty"`
	ContactEmail  string

	Reseller      string
	HomePartition string
	Shell         string
	Package       string `json:"package"`
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

	BandwidthUsed      *NumericLimit `json:"bandwidthused,omitempty"`
	BandwidthLimit     *NumericLimit `json:"bandwidthlimit,omitempty"`
	DiskUsed           *NumericLimit `json:"diskused,omitempty"`
	DiskLimit          *NumericLimit `json:"disklimit,omitempty"`
	AlternateDiskLimit *NumericLimit `json:"diskquota,omitempty"`
	InodeUsed          *NumericLimit `json:"inodeused,omitempty"`
	InodeLimit         *NumericLimit `json:"inodequota,omitempty"`
}

func (a *WhmAPI) ListAccounts() ([]string, error) {
	var outputData struct {
		Data struct {
			AccountList []struct {
				Username string `json:"user"`
			} `json:"acct"`
		} `json:"data"`
	}

	queryParams := url.Values{}
	queryParams.Add("want", "user")

	err := a.Call(
		http.MethodGet,
		"listaccts",
		queryParams,
		&outputData,
	)

	if err != nil {
		return []string{}, err
	}

	userlist := []string{}
	for _, account := range outputData.Data.AccountList {
		userlist = append(userlist, account.Username)
	}
	return userlist, nil
}
