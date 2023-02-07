package dirgroupgeo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGeoGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"account_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"codes": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Set:      schema.HashString,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
