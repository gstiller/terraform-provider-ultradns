package dirgroup

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func geoIPGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"codes": &schema.Schema{
			Type:     schema.TypeSet,
			Set:      schema.HashString,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
