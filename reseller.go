package cpanel

import (
	"net/http"
	"net/url"
)

func (a *WhmAPI) ListAllResellerNames() ([]string, error) {
	var outputData struct {
		Data struct {
			Resellers []string `json:"reseller"`
		} `json:"data"`
	}
	err := a.Call(
		http.MethodGet,
		"listresellers",
		url.Values{},
		&outputData,
	)

	if err != nil {
		return []string{}, err
	}

	return outputData.Data.Resellers, nil
}

type Reseller struct {
	User     string    `json:"user"`
	Accounts []Account `json:"acct"`

	BandwidthUsed        *IntLimit `json:"totalbwused,omitempty"`
	BandwidthAlloc       *IntLimit `json:"totalbwalloc,omitempty"`
	BandwidthLimit       *IntLimit `json:"bandwidthlimit,omitempty"`
	BandwidthOverSelling *CpBool   `json:"bwoverselling,omitempty"`

	DiskUsed        *FloatLimit `json:"diskused,omitempty"`
	DiskAlloc       *IntLimit   `json:"totaldiskalloc,omitempty"`
	DiskLimit       *IntLimit   `json:"diskquota,omitempty"`
	DiskOverselling *CpBool     `json:"diskoverselling,omitempty"`
}

func (a *WhmAPI) ResellerUsers(reseller string) (Reseller, error) {
	var outputData struct {
		Data struct {
			ResellerUser Reseller `json:"reseller"`
		} `json:"data"`
	}
	queryParams := url.Values{}
	queryParams.Add("user", reseller)
	queryParams.Add("filter_deleted", "1")

	err := a.Call(
		http.MethodGet,
		"resellerstats",
		queryParams,
		&outputData,
	)

	if err != nil {
		return Reseller{}, err
	}

	for id := range outputData.Data.ResellerUser.Accounts {
		outputData.Data.ResellerUser.Accounts[id].DiskLimit =
			outputData.Data.ResellerUser.Accounts[id].AlternateDiskLimit
	}

	return outputData.Data.ResellerUser, nil
}
