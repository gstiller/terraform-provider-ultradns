package dirgroupgeo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/helper"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/dirgroup/geo"
)

func ResourceGeoGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeoGroupCreate,
		ReadContext:   resourceGeoGroupRead,
		UpdateContext: resourceGeoGroupUpdate,
		DeleteContext: resourceGeoGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceGeoGroupSchema(),
	}
}

func resourceGeoGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	geoGroupData := newGeoGroup(rd)

	res, err := services.DirGroupGeoService.CreateDirGroupGeo(geoGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	uri := res.Header.Get("Location")
	id := helper.GetGeoIdFromURI(uri)
	rd.SetId(id)

	return resourceGeoGroupRead(ctx, rd, meta)
}

func resourceGeoGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var accountName string

	services := meta.(*service.Service)
	geoGroupData := newGeoGroup(rd)

	//	if val, ok := rd.GetOk("account_name"); ok {
	//		accountName = val.(string)
	//	}

	_, geoGroup, err := services.DirGroupGeoService.ReadDirGroupGeo(geoGroupData)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(geoGroup.Name)
	rd.Set("name", geoGroup.Name)
	rd.Set("account_name", accountName)
	rd.Set("description", geoGroup.Description)
	rd.Set("codes", geoGroup.Codes)

	return diags
}

func resourceGeoGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//var geoGroupName string

	services := meta.(*service.Service)
	geoGroupData := newGeoGroup(rd)

	//if val, ok := rd.GetOk("name"); ok {
	//	geoGroupName = val.(string)
	//}

	_, err := services.DirGroupGeoService.UpdateDirGroupGeo(geoGroupData)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGeoGroupRead(ctx, rd, meta)
}

func resourceGeoGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	//var accountName string

	services := meta.(*service.Service)
	geoGroupData := newGeoGroup(rd)
	//geoGroupName := rd.Id()

	//if val, ok := rd.GetOk("account_name"); ok {
	//	accountName = val.(string)
	//}
	_, err := services.DirGroupGeoService.DeleteDirGroupGeo(geoGroupData)

	if err != nil {
		rd.SetId("")

		return diag.FromErr(err)
	}

	rd.SetId("")

	return diags
}

func newGeoGroup(rd *schema.ResourceData) *geo.DirGroupGeo {
	geoData := &geo.DirGroupGeo{}

	if val, ok := rd.GetOk("name"); ok {
		geoData.Name = val.(string)
	}
	if val, ok := rd.GetOk("account_name"); ok {
		geoData.Account = val.(string)
	}
	if val, ok := rd.GetOk("description"); ok {
		geoData.Description = val.(string)
	}
	if val, ok := rd.GetOk("codes"); ok {
		log.Printf("val: %v, %T", val, val)
		geoCodesData := val.(*schema.Set).List()
		geoData.Codes = make([]string, len(geoCodesData))
		for i, geoCode := range geoCodesData {
			geoData.Codes[i] = geoCode.(string)
		}
	}

	return geoData
}
