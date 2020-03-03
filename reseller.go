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
