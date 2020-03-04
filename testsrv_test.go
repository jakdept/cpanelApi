package cpanel

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
)

var cpanelTestSrv httptest.Server
var mux http.ServeMux

func FormatQueryParams(values url.Values) string {
	var mapArgs []string

	for key, array := range values {
		sort.StringSlice(array)
		mapArgs = append(mapArgs, key+"_"+strings.Join(array, "-"))
	}

	sort.StringSlice(mapArgs)
	return strings.Join(mapArgs, ".")
}
