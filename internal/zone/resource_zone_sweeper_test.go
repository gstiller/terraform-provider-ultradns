package zone_test

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ultradns/terraform-provider-ultradns/internal/acctest"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/helper"
)

const testZoneSweeperPrefix = "terraform-plugin-acc-test-"

func TestMain(m *testing.M) {
	acctest.TestAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(make(map[string]interface{})))
	resource.TestMain(m)
}

func init() {
	resource.AddTestSweepers("ultradns_zone", &resource.Sweeper{
		Name: "ultradns_zone",
		F:    testAccZoneSweeper,
	})
}

func testAccZoneSweeper(r string) error {
	services := acctest.TestAccProvider.Meta().(*service.Service)
	query := testAccGetZoneSweeperQueryString()
	_, zoneList, err := services.ZoneService.ListZone(query)
	if err != nil {
		return err
	}

	for _, zone := range zoneList.Zones {
		if strings.HasPrefix(zone.Properties.Name, testZoneSweeperPrefix) {
			_, er := services.ZoneService.DeleteZone(zone.Properties.Name)

			if er != nil {
				log.Printf("error destroying %s during sweep: %s", zone.Properties.Name, er)
			}
		}
	}

	return nil
}

func testAccGetZoneSweeperQueryString() *helper.QueryInfo {
	return &helper.QueryInfo{
		Limit: 1000,
		Query: "name:" + testZoneSweeperPrefix,
	}
}
