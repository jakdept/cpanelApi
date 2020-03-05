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

	BandwidthUsed        *NumericLimit `json:"totalbwused,omitempty"`
	BandwidthAlloc       *NumericLimit `json:"totalbwalloc,omitempty"`
	BandwidthLimit       *NumericLimit `json:"bandwidthlimit,omitempty"`
	BandwidthOverSelling *NumericLimit `json:"bwoverselling,omitempty"`

	DiskUsed        *NumericLimit `json:"diskused,omitempty"`
	DiskAlloc       *NumericLimit `json:"totaldiskalloc,omitempty"`
	DiskLimit       *NumericLimit `json:"diskquota,omitempty"`
	DiskOverselling *NumericLimit `json:"diskoverselling,omitempty"`
}

func (a *WhmAPI) ListResellerUsers(reseller string) (Reseller, error) {
	var outputData struct {
		Data Reseller `json:"data"`
	}
	queryParams := url.Values{}
	queryParams.Add("user", "reseller")
	queryParams.Add("filter_deleted", "1")

	err := a.Call(
		http.MethodGet,
		"listresellers",
		queryParams,
		&outputData,
	)

	if err != nil {
		return Reseller{}, err
	}

	for id := range outputData.Data.Accounts {
		outputData.Data.Accounts[id].DiskLimit = outputData.Data.Accounts[id].AlternateDiskLimit
	}

	return outputData.Data, nil
}
