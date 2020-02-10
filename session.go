package cpapi

import "net/url"

type CpanelApi struct {
	url    url.URL
	domain string
}

func NewLocalApi() CpanelApi {
	session := CpanelApi{
		url: net.URL{
			Host: "https://hostname.example.com:2087"
			Path: "/cpsess##########/json-api/accountsummary?api.version=1&user=username"


	}

}
