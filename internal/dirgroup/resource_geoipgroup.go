package dirgroup

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/helper"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/dirgroup/geoip"
)

func ResourceGeoIPGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeoIPGroupCreate,
		ReadContext:   resourceGeoIPGroupRead,
		UpdateContext: resourceGeoIPGroupUpdate,
		DeleteContext: resourceGeoIPGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceGeoIPGroupSchema(),
	}
}

func resourceGeoIPGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	geoIPGroupData := newGeoIPGroup(rd)

	res, err := services.GroupGeoIPService.CreateDirGroupGeoIP(geoIPGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	uri := res.Header.Get("Location")
	id := helper.GetGeoIPIdFromURI(uri)
	rd.SetId(id)

	return resourceGeoIPGroupRead(ctx, rd, meta)
}

func resourceGeoIPGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var accountName string
	var geoIPGroupName string
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

func resourceGeoIPGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var geoIPGroupName string
	services := meta.(*service.Service)

	if val, ok := rd.GetOk("name"); ok {
		geoIPGroupName = val.(string)
	}

	geoIPGroupData := newGeoIPGroup(rd)

	_, err := services.GroupGeoIPService.UpdateDirGroupGeoIP(geoIPGroupName, geoIPGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGeoIPGroupRead(ctx, rd, meta)
}

func resourceGeoIPGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var accountName string

	services := meta.(*service.Service)
	geoIPGroupName := rd.Id()

	if val, ok := rd.GetOk("account_name"); ok {
		accountName = val.(string)
	}
	_, err := services.GroupGeoIPService.DeleteDirGroupGeoIP(geoIPGroupName, accountName)

	if err != nil {
		rd.SetId("")

		return diag.FromErr(err)
	}

	rd.SetId("")

	return diags
}

func newGeoIPGroup(rd *schema.ResourceData) *geoip.DirGroupGeoIP {
	geoIPData := &geoip.DirGroupGeoIP{}

	if val, ok := rd.GetOk("name"); ok {
		geoIPData.Name = val.(string)
	}
	if val, ok := rd.GetOk("account_name"); ok {
		geoIPData.Account = val.(string)
	}
	if val, ok := rd.GetOk("description"); ok {
		geoIPData.Description = val.(string)
	}
	if val, ok := rd.GetOk("codes"); ok {
		log.Printf("val: %v, %T", val, val)
		geoCodesData := val.(*schema.Set).List()
		geoIPData.Codes = make([]string, len(geoCodesData))
		for i, geoCode := range geoCodesData {
			geoIPData.Codes[i] = geoCode.(string)
		}
	}

	return geoIPData
}
