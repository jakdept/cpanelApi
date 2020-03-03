package cpanel

import (
	"net/http"
	"net/url"
)

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
