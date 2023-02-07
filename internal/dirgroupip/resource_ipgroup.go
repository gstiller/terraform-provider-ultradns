package dirgroupip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/helper"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/dirgroup/ip"
)

func ResourceIPGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGroupCreate,
		ReadContext:   resourceIPGroupRead,
		UpdateContext: resourceIPGroupUpdate,
		DeleteContext: resourceIPGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceIPGroupSchema(),
	}
}

func resourceIPGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	ipGroupData := newIPGroup(rd)

	res, err := services.DirGroupIPService.CreateDirGroupIP(ipGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	uri := res.Header.Get("Location")
	id := helper.GetIPIdFromURI(uri)
	rd.SetId(id)

	return resourceIPGroupRead(ctx, rd, meta)
}

func resourceIPGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var accountName string

	services := meta.(*service.Service)
	ipGroupData := newIPGroup(rd)

	if val, ok := rd.GetOk("account_name"); ok {
		accountName = val.(string)
	}
	_, ipGroup, err := services.DirGroupIPService.ReadDirGroupIP(ipGroupData)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(ipGroup.Name)
	rd.Set("name", ipGroup.Name)
	rd.Set("account_name", accountName)
	rd.Set("description", ipGroup.Description)
	rd.Set("ips", ipGroup.IPs)

	return diags
}

func resourceIPGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	ipGroupData := newIPGroup(rd)

	_, err := services.DirGroupIPService.UpdateDirGroupIP(ipGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIPGroupRead(ctx, rd, meta)
}

func resourceIPGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	services := meta.(*service.Service)
	ipGroupData := newIPGroup(rd)
	_, err := services.DirGroupIPService.DeleteDirGroupIP(ipGroupData)

	if err != nil {
		rd.SetId("")

		return diag.FromErr(err)
	}

	rd.SetId("")

	return diags
}

func newIPGroup(rd *schema.ResourceData) *ip.DirGroupIP {
	ipData := &ip.DirGroupIP{}

	if val, ok := rd.GetOk("name"); ok {
		ipData.Name = val.(string)
	}
	if val, ok := rd.GetOk("account_name"); ok {
		ipData.Account = val.(string)
	}
	if val, ok := rd.GetOk("description"); ok {
		ipData.Description = val.(string)
	}
	if val, ok := rd.GetOk("ip"); ok {
		sourceIPAddressDataList := val.(*schema.Set).List()
		ipData.IPs = getSourceIPAddressList(sourceIPAddressDataList)
	}

	return ipData
}

func getSourceIPAddressList(sourceIPAddressDataList []interface{}) []*ip.IPAddress {
	sourceIPAddressList := make([]*ip.IPAddress, len(sourceIPAddressDataList))

	for i, d := range sourceIPAddressDataList {
		sourceIPAddressData := d.(map[string]interface{})
		sourceIPAddressList[i] = getSourceIPAddress(sourceIPAddressData)
	}

	return sourceIPAddressList
}

func getSourceIPAddress(sourceIPAddressData map[string]interface{}) *ip.IPAddress {
	sourceIPAddress := &ip.IPAddress{}

	if val, ok := sourceIPAddressData["start"]; ok {
		sourceIPAddress.Start = val.(string)
	}

	if val, ok := sourceIPAddressData["end"]; ok {
		sourceIPAddress.End = val.(string)
	}

	if val, ok := sourceIPAddressData["cidr"]; ok {
		sourceIPAddress.Cidr = val.(string)
	}

	if val, ok := sourceIPAddressData["address"]; ok {
		sourceIPAddress.Address = val.(string)
	}

	return sourceIPAddress
}
