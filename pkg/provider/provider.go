package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/Verlet64/terraform-provider-insightops/pkg/insightops"
	"github.com/Verlet64/terraform-provider-insightops/pkg/resources"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"insightops_saved_query": resources.ResourceInsightOpsSavedQuery(),
		},
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: configureInsightopsProvider,
	}
}

func configureInsightopsProvider(d *schema.ResourceData) (interface{}, error) {
	return insightops.NewClient(d.Get("api_key").(string), d.Get("endpoint").(string)), nil
}
