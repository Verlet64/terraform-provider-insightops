package savedqueries_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Verlet64/terraform-provider-insightops/pkg/insightops/savedqueries"
)

func getMockUpdateKeys() map[string]string {
	keys := make(map[string]string)

	keys["valid-api-key"] = "auth"

	keys["updateable-query-id"] = "id"
	keys["not-found-query-id"] = "id2"
	keys["internal-server-err-query-id"] = "id3"
	keys["bad-media-err-query-id"] = "id4"
	keys["bad-method-err-query-id"] = "id5"

	return keys
}

func getTestUpdateServer(keys map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != keys["valid-api-key"] {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		split := strings.Split(r.URL.Path, "/")
		id := split[len(split)-1]

		fmt.Printf("%v", id)

		if id == keys["internal-server-err-query-id"] {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if id == keys["bad-method-err-query-id"] {
			w.WriteHeader(http.StatusMethodNotAllowed)
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
		w.Write([]byte(`{"saved_query": {"name": "name", "leql": {"statement": "where(test)"}}}`))
		return
	}))

	return server
}

func TestUpdateSavedQueryBadAuth(t *testing.T) {
	keys := getMockUpdateKeys()
	server := getTestUpdateServer(keys)
	defer server.Close()

	expected := errors.New("You are not authorised to perform this action. [Status 403]")
	result, err := savedqueries.UpdateSavedQuery(server.URL, "", keys["updateable-query-id"], "name", "query")

	if result != nil {
		t.Errorf("Expected error, got: %v", result)
		return
	}

	if err == nil {
		t.Errorf("Expected: %v, got nil", expected)
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, err)
	}
}

func TestUpdateSavedQueryResourceNotFound(t *testing.T) {
	keys := getMockUpdateKeys()
	server := getTestUpdateServer(keys)
	defer server.Close()

	expected := errors.New("Unable to locate saved query. [Status 404]")
	result, err := savedqueries.UpdateSavedQuery(server.URL, keys["valid-api-key"], keys["not-found-query-id"], "name", "query")

	if result != nil {
		t.Errorf("Expected error, got: %v", result)
		return
	}

	if err == nil {
		t.Errorf("Expected: %v, got nil", expected)
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, err)
	}

}

func TestUpdateSavedQueryServerDown(t *testing.T) {
	keys := getMockUpdateKeys()
	server := getTestUpdateServer(keys)
	defer server.Close()

	expected := errors.New("Unable to modify the query at this time. [Status 500]")
	result, err := savedqueries.UpdateSavedQuery(server.URL, keys["valid-api-key"], keys["internal-server-err-query-id"], "name", "query")

	if result != nil {
		t.Errorf("Expected error, got: %v", result)
		return
	}

	if err == nil {
		t.Errorf("Expected: %v, got nil", expected)
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, err)
	}
}

func TestUpdateSavedQueryBadQuery(t *testing.T) {
	keys := getMockUpdateKeys()
	server := getTestUpdateServer(keys)
	defer server.Close()

	expected := errors.New("Unable to modify the query at this time. Please raise an issue. [Status 415]")
	result, err := savedqueries.UpdateSavedQuery(server.URL, keys["valid-api-key"], keys["bad-media-err-query-id"], "name", "query")

	if result != nil {
		t.Errorf("Expected error, got: %v", result)
		return
	}

	if err == nil {
		t.Errorf("Expected: %v, got nil", expected)
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, err)
	}
}

func TestUpdateSavedQueryBadMethod(t *testing.T) {
	keys := getMockUpdateKeys()
	server := getTestUpdateServer(keys)
	defer server.Close()

	expected := errors.New("Unable to modify the query at this time. Please raise an issue. [Status 405]")
	result, err := savedqueries.UpdateSavedQuery(server.URL, keys["valid-api-key"], keys["bad-method-err-query-id"], "name", "query")

	if result != nil {
		t.Errorf("Expected error, got: %v", result)
		return
	}

	if err == nil {
		t.Errorf("Expected: %v, got nil", expected)
	}

	if err.Error() != expected.Error() {
		t.Errorf("Expected: %v, got: %v", expected, err)
	}
}
