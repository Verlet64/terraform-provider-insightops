package savedqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CreateSavedQuery creates a saved query in InsightOps
//
// It doesn't handle any form of parsing/validation of
// the query itself, so the query may be invalid.
func CreateSavedQuery(uri string, apikey string, name string, query string) (*SavedQueryResponse, error) {
	type During struct {
		To        interface{} `json:"to"`
		From      interface{} `json:"from"`
		TimeRange interface{} `json:"time_range"`
	}

	type Leql struct {
		During    During `json:"during"`
		Statement string `json:"statement"`
	}

	type SavedQuery struct {
		Logs []interface{} `json:"logs"`
		Leql Leql          `json:"leql"`
		Name string        `json:"name"`
	}

	type SavedQueryRequest struct {
		SavedQuery SavedQuery `json:"saved_query"`
	}

	body := SavedQueryRequest{SavedQuery: SavedQuery{Name: name, Leql: Leql{Statement: query}}}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("POST", uri, buf)
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
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("You are not authorised to perform this action. [Status %v]", res.StatusCode)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to create a saved query at this time. [Status %v]", res.StatusCode)
	}

	var response SavedQueryResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
