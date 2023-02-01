package dirgroup

import (
	"context"
	// "errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
)

func ResourceGeoIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: geoIPCreate,
		ReadContext:   geoIPRead,
		UpdateContext: geoIPUpdate,
		DeleteContext: geoIPDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// DiffSupressFunc?
				// ValidateDiagFunc?
			},
			"description": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
			},
			"codes": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func geoIPCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func geoIPRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var groupName string
	services := meta.(*service.Service)

	if val, ok := rd.GetOk("name"); ok {
		groupName = val.(string)
	}

	log.Printf("THE DDATAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", groupName)
	_, geoIPGroup, err := services.GroupGeoIPService.Read(groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(geoIPGroup.Name)
	rd.Set("name", geoIPGroup.Name)
	rd.Set("description", geoIPGroup.Description)
	rd.Set("codes", geoIPGroup.Codes)
	log.Printf("Look, a group?", geoIPGroup)
	return diags
}

func geoIPUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func geoIPDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
