package insightops

import "example.com/terraform-provider-insightops/pkg/insightops/savedqueries"

type insightopsclientiface interface {
	FetchSavedQuery(id string) (*savedqueries.SavedQueryResponse, error)
	CreateSavedQuery(name string, query string) (*savedqueries.SavedQueryResponse, error)
	UpdateSavedQuery(id string, name string, query string) (*savedqueries.SavedQueryResponse, error)
	DeleteSavedQuery(id string) error
}

type Client struct {
	APIKey string
	URI    string
}

func NewClient(key string, endpoint string) insightopsclientiface {
	return &Client{
		APIKey: key,
		URI:    endpoint,
	}
}

func (c *Client) FetchSavedQuery(id string) (*savedqueries.SavedQueryResponse, error) {
	res, err := savedqueries.FetchSavedQuery(c.URI, c.APIKey, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteSavedQuery(id string) error {
	err := savedqueries.DeleteSavedQuery(c.URI, c.APIKey, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateSavedQuery(name string, query string) (*savedqueries.SavedQueryResponse, error) {
	res, err := savedqueries.CreateSavedQuery(c.URI, c.APIKey, name, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) UpdateSavedQuery(id string, name string, query string) (*savedqueries.SavedQueryResponse, error) {
	res, err := savedqueries.UpdateSavedQuery(c.APIKey, id, name, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}
