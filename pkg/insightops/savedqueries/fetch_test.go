package savedqueries_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Verlet64/terraform-provider-insightops/pkg/insightops/savedqueries"
)

func getMockFetchKeys() map[string]string {
	keys := make(map[string]string)

	keys["valid-api-key"] = "test-key"
	keys["valid-id"] = "id"
	keys["invalid-media-id"] = "id-2"
	keys["invalid-method-id"] = "id-3"

	return keys
}

func getTestFetchServer(keys map[string]string, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != keys["valid-api-key"] {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}

		split := strings.Split(r.URL.Path, "/")
		id := split[len(split)-1]

		if id == keys["invalid-media-id"] {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		if id == keys["invalid-method-id"] {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if id != keys["valid-id"] {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))

		return
	}))

	return server
}

func TestFetchQuery(t *testing.T) {
	keys := getMockFetchKeys()
	server := getTestFetchServer(keys, `{"saved_query": {"name": "name", "leql": {"statement": "where(test)"}}}`)
	defer server.Close()

	result, err := savedqueries.FetchSavedQuery(server.URL, keys["valid-api-key"], keys["valid-id"])
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	got := result.SavedQuery.Name
	expected := "name"

	if got != expected {
		t.Errorf("Expected query name to be %v, got %v", expected, got)
	}

	got = result.SavedQuery.Leql.Statement
	expected = "where(test)"

	if got != expected {
		t.Errorf("Expected query statement to be %v, got %v", expected, got)
	}

}

func TestFetchQueryInvalidMethod(t *testing.T) {
	keys := getMockFetchKeys()
	server := getTestFetchServer(keys, `{`)
	defer server.Close()

	expected := "Unable to fetch the saved query at this time. Please raise an issue. [Status 405]"

	result, err := savedqueries.FetchSavedQuery(server.URL, keys["valid-api-key"], keys["invalid-method-id"])
	if result != nil {
		t.Errorf("Expected no result, got: %v", result)
	}

	if err.Error() != expected {
		t.Errorf("Expected %v, got: %v", expected, err.Error())
		return
	}
}

func TestFetchQueryInvalidFormat(t *testing.T) {
	keys := getMockFetchKeys()
	server := getTestFetchServer(keys, `{`)
	defer server.Close()

	expected := "Unable to fetch the saved query at this time. Please raise an issue. [Status 415]"

	result, err := savedqueries.FetchSavedQuery(server.URL, keys["valid-api-key"], keys["invalid-media-id"])
	if result != nil {
		t.Errorf("Expected no result, got: %v", result)
	}

	if err.Error() != expected {
		t.Errorf("Expected %v, got: %v", expected, err.Error())
		return
	}
}

func TestFetchQueryNonExistentId(t *testing.T) {
	keys := getMockFetchKeys()
	server := getTestFetchServer(keys, `{"saved_query": {"name": "name", "leql": {"statement": "where(test)"}}}`)
	defer server.Close()

	expected := "Unable to locate saved query. [Status 404]"

	result, err := savedqueries.FetchSavedQuery(server.URL, keys["valid-api-key"], "random-id")
	if result != nil {
		t.Errorf("Expected no result, got: %v", result)
	}

	if err == nil {
		t.Errorf("Expected an error, got nothing")
	}

	if err.Error() != expected {
		t.Errorf("Expected %v, got %v", expected, err.Error())
	}
}
