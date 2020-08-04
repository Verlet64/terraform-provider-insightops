package savedqueries_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Verlet64/terraform-provider-insightops/pkg/insightops/savedqueries"
)

func getMockDeleteKeys() map[string]string {
	keys := make(map[string]string)

	keys["valid-api-key"] = "auth"

	keys["deleteable-query-id"] = "id"
	keys["not-found-query-id"] = "id2"
	keys["internal-server-err-query-id"] = "id3"
	keys["bad-media-err-query-id"] = "id4"

	return keys
}

func getTestDeleteServer(keys map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != keys["valid-api-key"] {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		split := strings.Split(r.URL.Path, "/")
		id := split[len(split)-1]

		if id == keys["internal-server-err-query-id"] {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if id == keys["bad-media-err-query-id"] {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		if id == keys["not-found-query-id"] {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}))
}

func TestDeleteSavedQuery(t *testing.T) {
	keys := getMockDeleteKeys()
	server := getTestDeleteServer(keys)
	defer server.Close()

	got := savedqueries.DeleteSavedQuery(server.URL, keys["valid-api-key"], keys["deleteable-query-id"])

	if got != nil {
		t.Errorf("Expected: %v, Got: %v", nil, got)
	}
}

func TestDeleteSavedQueryBadAuth(t *testing.T) {
	keys := getMockDeleteKeys()
	server := getTestDeleteServer(keys)
	defer server.Close()

	expected := errors.New("You are not authorised to perform this action. [Status 403]")
	got := savedqueries.DeleteSavedQuery(server.URL, "", keys["deleteable-query-id"])

	if got == nil {
		t.Errorf("Expected: %v, got: %v", expected, nil)
		return
	}

	if got.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
}

func TestDeleteSavedQueryResourceNotFound(t *testing.T) {
	keys := getMockDeleteKeys()
	server := getTestDeleteServer(keys)
	defer server.Close()

	got := savedqueries.DeleteSavedQuery(server.URL, keys["valid-api-key"], keys["not-found-query-id"])

	if got != nil {
		t.Errorf("Expected: %v, got: %v", nil, got)
	}
}

func TestDeleteSavedQueryServerDown(t *testing.T) {
	keys := getMockDeleteKeys()
	server := getTestDeleteServer(keys)
	defer server.Close()

	expected := errors.New("Unable to modify the query at this time. [Status 500]")
	got := savedqueries.DeleteSavedQuery(server.URL, keys["valid-api-key"], keys["internal-server-err-query-id"])

	if got == nil {
		t.Errorf("Expected: %v, got: %v", expected, nil)
		return
	}

	if got.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
}

func TestDeleteSavedQueryBadQuery(t *testing.T) {
	keys := getMockDeleteKeys()
	server := getTestDeleteServer(keys)
	defer server.Close()

	expected := errors.New("Unable to modify the query at this time. Please raise an issue. [Status 415]")
	got := savedqueries.DeleteSavedQuery(server.URL, keys["valid-api-key"], keys["bad-media-err-query-id"])

	if got == nil {
		t.Errorf("Expected: %v, got: %v", expected, nil)
		return
	}

	if got.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
}
