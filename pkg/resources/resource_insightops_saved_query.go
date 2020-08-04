package resources

import (
	"github.com/Verlet64/terraform-provider-insightops/pkg/insightops"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInsightOpsSavedQueryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*insightops.Client)

	name := d.Get("name").(string)
	query := d.Get("query").(string)

	res, err := client.CreateSavedQuery(name, query)

	if err != nil {
		return err
	}

	d.SetId(res.SavedQuery.ID)

	return resourceInsightOpsSavedQueryRead(d, m)
}

func resourceInsightOpsSavedQueryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*insightops.Client)

	id := d.Id()

	res, err := client.FetchSavedQuery(id)
	if err != nil && err.Error() == "not found" {
		d.SetId("")

		return nil
	}

	d.Set("name", res.SavedQuery.Name)
	d.Set("query", res.SavedQuery.Leql.Statement)

	return err
}

func resourceInsightOpsSavedQueryUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*insightops.Client)

	id := d.Id()

	name := d.Get("name").(string)
	query := d.Get("query").(string)

	res, err := client.UpdateSavedQuery(id, name, query)
	if err != nil {
		return err
	}

	d.Set("name", res.SavedQuery.Name)
	d.Set("query", res.SavedQuery.Leql.Statement)

	return nil
}

func resourceInsightOpsSavedQueryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*insightops.Client)

	id := d.Id()

	err := client.DeleteSavedQuery(id)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func ResourceInsightOpsSavedQuery() *schema.Resource {
	return &schema.Resource{
		Create: resourceInsightOpsSavedQueryCreate,
		Read:   resourceInsightOpsSavedQueryRead,
		Update: resourceInsightOpsSavedQueryUpdate,
		Delete: resourceInsightOpsSavedQueryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_selection": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"time_range": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
