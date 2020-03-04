package cpanel

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var cpanelTestSrv httptest.Server
var mux http.ServeMux

var testAuthToken = "/cpsess8675309"

// This handler checks for the API version and returns a 404 if not found.
func StripApiVersion(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query.Get("api.version")) != 1 ||
			r.URL.Query.Get("api.version")[0] != "1" {
			http.NotFound(w, r)
			return
		}
		r.URL.Query.Delete("api.version")
		h.ServeHTTP(w, r)
	})
}

func StripToken(authHandler http.Handler, unauthHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, testAuthToken); len(p) < len(r.URL.Path) {
			r2 := new(Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = p
			authHandler.ServeHTTP(w, r2)
		} else {
			unauthHandler.ServeHTTP(w, r)
		}
	})
}

func RespondWithFile(pathPrefix string, notFoundCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(filepath.Join(
			strings.TrimPrefix(pathPrefix, "/"),
			r.URL.Path,
			r.URL.Query().Encode,
		))
		if err != nil && os.IsNotExist(err) {
			http.Error(w, notFoundCode, "")
		} else if err != nil {
			http.Error(w, http.StatusInternalServerError, "internal problem")
		} else {
			_, err = io.Copy(w, f)
			if err != nil {
				http.Error(w, notFoundCode, "")
			}
		}
	})
}
