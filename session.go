package cpanel

import (
	"fmt"
	"net"
	"net/url"
)

type WhmAPI struct {
	hostname string
	token    string
}

// TODO figure out how to standardize ssh.Session and exec.Command

func (a *WhmAPI) GenerateURL(endpoint string) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     a.hostname + ":2087",
		Path:     fmt.Sprintf("/%s/json-api/%s", a.token, endpoint),
		RawQuery: "?api.version=1",
	}
}

type cPanelAccount struct {
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

func (a *WhmAPI) ListAccounts() {

}
