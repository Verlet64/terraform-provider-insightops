package savedqueries_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/terraform-provider-insightops/pkg/insightops/savedqueries"
)

func getMockCreateKeys() map[string]string {
	keys := make(map[string]string)

	keys["test-api-key"] = "test-key"

	return keys
}

func getTestCreateServer(keys map[string]string, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("x-api-key") != keys["test-api-key"] {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))

		return
	}))

	return server
}

func TestCreateSavedQuery(t *testing.T) {
	keys := getMockCreateKeys()
	server := getTestCreateServer(keys, `{"saved_query": {"name": "name", "leql": {"statement": "where(test)"}}}`)
	defer server.Close()

	result, err := savedqueries.CreateSavedQuery(server.URL, keys["test-api-key"], "name", "where(test)")
	if err != nil {
		t.Errorf("Failed to query API with error: %v", err)
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

func TestCreateSavedQueryBadAuth(t *testing.T) {
	keys := getMockCreateKeys()
	server := getTestCreateServer(keys, `{"saved_query": {"name": "name", "leql": {"statement": "where(test)"}}}`)
	defer server.Close()

	result, err := savedqueries.CreateSavedQuery(server.URL, "", "name", "where(test)")
	if result != nil {
		t.Errorf("Expected an authorisation error")
	}

	expected := "You are not authorised to perform this action. [Status 403]"
	if err.Error() != expected {
		t.Errorf("Expected: %v, Got: %v", expected, err.Error())
	}
}
