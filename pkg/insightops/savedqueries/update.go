package savedqueries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// UpdateSavedQuery updates a query in Insightops.
// Currently it can only update the name and query of a saved query.
func UpdateSavedQuery(uri string, apikey string, id string, name string, query string) (*SavedQueryResponse, error) {
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

	queryURI := strings.Join([]string{uri, id}, "/")

	body := SavedQueryRequest{SavedQuery: SavedQuery{Name: name, Leql: Leql{Statement: query}}}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("PATCH", queryURI, buf)
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

	if res.StatusCode >= 400 {
		switch res.StatusCode {
		case http.StatusForbidden:
			return nil, fmt.Errorf("You are not authorised to perform this action. [Status %v]", http.StatusForbidden)
		case http.StatusNotFound:
			return nil, fmt.Errorf("Unable to locate saved query. [Status %v]", http.StatusNotFound)
		case http.StatusInternalServerError:
			return nil, fmt.Errorf("Unable to modify the query at this time. [Status %d]", http.StatusInternalServerError)
		case http.StatusMethodNotAllowed:
			fallthrough
		case http.StatusUnsupportedMediaType:
			fallthrough
		default:
			return nil, fmt.Errorf("Unable to modify the query at this time. Please raise an issue. [Status %d]", res.StatusCode)
		}
	}

	defer res.Body.Close()

	var response SavedQueryResponse

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
