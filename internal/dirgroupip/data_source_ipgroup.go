package dirgroupip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
)

func DataSourceIPGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPGroupRead,

		Schema: dataSourceIPSchema(),
	}
}

func dataSourceIPGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	services := meta.(*service.Service)
	ipGroupData := newIPGroup(rd)

	_, ipGroup, err := services.DirGroupIPService.ReadDirGroupIP(ipGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(ipGroup.Name)
	rd.Set("name", ipGroup.Name)
	rd.Set("account_name", ipGroup.Account)
	rd.Set("description", ipGroup.Description)
	rd.Set("ip", ipGroup.IPs)
	return diags

}
