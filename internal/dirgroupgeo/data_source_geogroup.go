package dirgroupgeo

import (
	"context"
	//	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	// "github.com/ultradns/ultradns-go-sdk/pkg/dirgroup/geo"
)

func DataSourceGeoGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeoGroupRead,

		Schema: dataSourceGeoSchema(),
	}
}

func dataSourceGeoGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	services := meta.(*service.Service)
	geoGroupData := newGeoGroup(rd)

	//	if val, ok := rd.GetOk("account_name"); ok {
	//		accountName = val.(string)
	//	}
	//	if val, ok := rd.GetOk("name"); ok {
	//		geoGroupName = val.(string)
	//	}

	_, geoGroup, err := services.DirGroupGeoService.ReadDirGroupGeo(geoGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(geoGroup.Name)
	rd.Set("name", geoGroup.Name)
	rd.Set("account_name", geoGroup.Account)
	rd.Set("description", geoGroup.Description)
	rd.Set("codes", geoGroup.Codes)
	return diags

}

//	func listGeoGroup(ctx context.Context, rd *schema.ResourceData) diag.Diagnostics {
//		services := meta.(*service.Service)
//		geoGroup := &geoip.DirGroupGeo{}
//
//		if val, ok := rd.GetOk("account_name"); ok {
//			geoGroup.Account = val.(string)
//		}
//	}
