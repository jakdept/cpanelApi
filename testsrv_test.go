package cpanel

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
)

var testSrv *httptest.Server
var testAuthToken = "/cpsess8675309"
var testWhmApi WhmAPI

func TestMain(m *testing.M) {
	unauthHandler := RespondWithFile(
		"testdata/fixtures/testSrv/create_user_session",
		http.StatusUnauthorized,
	)
	unauthHandler = http.StripPrefix("/create_user_session", unauthHandler)
	unauthHandler = http.StripPrefix("/json-api/", unauthHandler)

	authHandler := RespondWithFile(
		"testdata/fixtures/testSrv",
		http.StatusNotFound,
	)
	authHandler = http.StripPrefix("/json-api/", authHandler)

	handler := StripToken(authHandler, unauthHandler)
	handler = StripApiVersion(handler)
	testSrv = httptest.NewTLSServer(handler)
	srvHost := strings.TrimPrefix(testSrv.URL, "https://")

	testWhmApi = WhmAPI{
		hostname: &srvHost,
		token:    &testAuthToken,
		client:   testSrv.Client(),
	}

	os.Exit(m.Run())
}

// This handler checks for the API version and returns a 404 if not found.
func StripApiVersion(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api.version") == "1" {
			params := r.URL.Query()
			params.Del("api.version")
			r.URL.RawQuery = params.Encode()
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

func StripToken(authHandler http.Handler, unauthHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, testAuthToken); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
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
		filename := r.URL.Query().Encode()
		if filename == "" {
			filename = "empty"
		}
		filename += ".json"
		filename = filepath.Join(
			pathPrefix,
			r.URL.Path,
			filename,
		)
		color.Cyan("opening %s\n", filename)

		f, err := os.Open(filename)
		if err != nil && os.IsNotExist(err) {
			http.Error(w, "", notFoundCode)
		} else if err != nil {
			http.Error(w, "internal problem", http.StatusInternalServerError)
		} else {
			_, err = io.Copy(w, f)
			if err != nil {
				http.Error(w, "", notFoundCode)
			}
		}
	})
}
