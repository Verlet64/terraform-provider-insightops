package savedqueries

// SavedQueryResponse is the response format for REST operation
// relating to a saved query.
type SavedQueryResponse struct {
	SavedQuery struct {
		ID   string        `json:"id"`
		Logs []interface{} `json:"logs"`
		Leql struct {
			During struct {
				To        interface{} `json:"to"`
				From      interface{} `json:"from"`
				TimeRange interface{} `json:"time_range"`
			} `json:"during"`
			Statement string `json:"statement"`
		} `json:"leql"`
		Name string `json:"name"`
	} `json:"saved_query"`
}
