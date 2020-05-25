package savedqueries

import (
	"net/http"
	"strings"
	"time"
)

// DetchSavedQuery
func DeleteSavedQuery(apikey string, id string) error {
	queryURI := strings.Join([]string{SavedQueriesBaseURI, id}, "/")

	req, err := http.NewRequest("DELETE", queryURI, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apikey)

	client := &http.Client{Timeout: time.Second * 10}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
