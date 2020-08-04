package savedqueries

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DeleteSavedQuery given an endpoint, api key and id for a saved query deletes
// the query in Rapid7
func DeleteSavedQuery(endpoint string, apikey string, id string) error {
	queryURI := strings.Join([]string{endpoint, id}, "/")

	req, err := http.NewRequest("DELETE", queryURI, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apikey)

	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		fallthrough
	case http.StatusNotFound:
		return nil
	case http.StatusUnsupportedMediaType:
		return fmt.Errorf("Unable to modify the query at this time. Please raise an issue. [Status %v]", http.StatusUnsupportedMediaType)
	case http.StatusForbidden:
		return fmt.Errorf("You are not authorised to perform this action. [Status %v]", res.StatusCode)
	default:
		return fmt.Errorf("Unable to modify the query at this time. [Status %v]", res.StatusCode)
	}

}
