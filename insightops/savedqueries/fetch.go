package savedqueries

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// FetchSavedQuery fetches a saved query from InsightOps
func FetchSavedQuery(apikey string, id string) (*SavedQueryResponse, error) {
	queryURI := strings.Join([]string{SavedQueriesBaseURI, id}, "/")

	req, err := http.NewRequest("GET", queryURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apikey)

	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	defer res.Body.Close()

	var response SavedQueryResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}