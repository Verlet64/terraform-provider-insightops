package main

import (
	"example.com/terraform-provider-insightops/insightops"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"insightops_saved_query": resourceInsightOpsSavedQuery(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: configureInsightopsProvider,
	}
}

func configureInsightopsProvider(d *schema.ResourceData) (interface{}, error) {
	return insightops.NewClient(d.Get("api_key").(string)), nil
}
