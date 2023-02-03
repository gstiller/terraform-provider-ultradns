package dirgroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	// "github.com/ultradns/ultradns-go-sdk/pkg/dirgroup/geoip"
)

func DataSourceGeoIPGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeoIPGroupRead,

		Schema: dataSourceGeoIPSchema(),
	}
}

func dataSourceGeoIPGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var geoIPGroupName string
	var accountName string
	services := meta.(*service.Service)

	if val, ok := rd.GetOk("account_name"); ok {
		accountName = val.(string)
	}
	if val, ok := rd.GetOk("name"); ok {
		geoIPGroupName = val.(string)
	}

	_, geoIPGroup, err := services.GroupGeoIPService.ReadDirGroupGeoIP(geoIPGroupName, accountName)

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(geoIPGroup.Name)
	rd.Set("name", geoIPGroup.Name)
	rd.Set("account_name", accountName)
	rd.Set("description", geoIPGroup.Description)
	rd.Set("codes", geoIPGroup.Codes)
	return diags

}

//func listGeoIPGroup(ctx context.Context, rd *schema.ResourceData) diag.Diagnostics {
//	services := meta.(*service.Service)
//	geoIPGroup := &geoip.DirGroupGeoIP{}
//
//	if val, ok := rd.GetOk("account_name"); ok {
//		geoIPGroup.Account = val.(string)
//	}
//}
