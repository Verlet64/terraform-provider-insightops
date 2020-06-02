package savedqueries

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// FetchSavedQuery fetches a saved query from InsightOps
func FetchSavedQuery(uri string, apikey string, id string) (*SavedQueryResponse, error) {
	queryURI := strings.Join([]string{uri, id}, "/")

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

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("Unable to locate saved query. [Status %v]", res.StatusCode)
	case http.StatusForbidden:
		return nil, fmt.Errorf("Unable to locate saved query. [Status %v]", res.StatusCode)
	case http.StatusMethodNotAllowed:
		fallthrough
	case http.StatusUnsupportedMediaType:
		return nil, fmt.Errorf("Unable to fetch the saved query at this time. Please raise an issue. [Status %v]", res.StatusCode)
	default:
	}

	defer res.Body.Close()

	var response SavedQueryResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errors.New("Unable to fetch the saved query at this time. Please raise an issue")
	}

	return &response, nil
}
