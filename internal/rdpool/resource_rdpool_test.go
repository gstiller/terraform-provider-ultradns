package rdpool_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ultradns/terraform-provider-ultradns/internal/acctest"
	"github.com/ultradns/terraform-provider-ultradns/internal/errors"
	"github.com/ultradns/terraform-provider-ultradns/internal/rrset"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/record/pool"
)

const zoneResourceName = "primary_rdpool"

func TestAccResourceRDPool(t *testing.T) {
	zoneName := acctest.GetRandomZoneName()
	ownerNameTypeA := tfacctest.RandString(3)
	ownerNameTypeAAAA := tfacctest.RandString(3)
	testCase := resource.TestCase{
		PreCheck:     func() { acctest.TestPreCheck(t) },
		Providers:    acctest.TestAccProviders,
		CheckDestroy: testAccCheckRDPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRDPoolA(zoneName, ownerNameTypeA),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDPoolExists("ultradns_rdpool.a"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "zone_name", zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "owner_name", ownerNameTypeA+"."+zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "record_type", "A"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "ttl", "120"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "record_data.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "order", "FIXED"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "description", "RD Pool Resource of Type A"),
				),
			},
			{
				Config: testAccResourceUpdateRDPoolA(zoneName, ownerNameTypeA),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDPoolExists("ultradns_rdpool.a"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "zone_name", zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "owner_name", ownerNameTypeA+"."+zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "record_type", "A"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "ttl", "150"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "record_data.0", "192.168.1.2"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "order", "RANDOM"),
					resource.TestCheckResourceAttr("ultradns_rdpool.a", "description", ownerNameTypeA+"."+zoneName),
				),
			},
			{
				ResourceName:      "ultradns_rdpool.a",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceRDPoolAAAA(zoneName, ownerNameTypeAAAA),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDPoolExists("ultradns_rdpool.aaaa"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "zone_name", zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "owner_name", ownerNameTypeAAAA+"."+zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "record_type", "AAAA"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "ttl", "120"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "record_data.0", "aaaa:bbbb:cccc:dddd:eeee:ffff:1111:2222"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "order", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "description", ownerNameTypeAAAA+"."+zoneName),
				),
			},
			{
				Config: testAccResourceUpdateRDPoolAAAA(zoneName, ownerNameTypeAAAA),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRDPoolExists("ultradns_rdpool.aaaa"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "zone_name", zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "owner_name", ownerNameTypeAAAA+"."+zoneName),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "record_type", "AAAA"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "ttl", "150"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "record_data.0", "aaaa:bbbb:cccc:dddd:eeee:ffff:1111:3333"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "order", "FIXED"),
					resource.TestCheckResourceAttr("ultradns_rdpool.aaaa", "description", "RD Pool Resource of Type AAAA"),
				),
			},
		},
	}
	resource.ParallelTest(t, testCase)
}

func testAccCheckRDPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return errors.ResourceNotFoundError(resourceName)
		}

		services := acctest.TestAccProvider.Meta().(*service.Service)
		rrSetKey := rrset.GetRRSetKeyFromID(rs.Primary.ID)
		rrSetKey.PType = pool.RD
		_, _, err := services.RecordService.Read(rrSetKey)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckRDPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ultradns_rdpool" {
			continue
		}

		services := acctest.TestAccProvider.Meta().(*service.Service)
		rrSetKey := rrset.GetRRSetKeyFromID(rs.Primary.ID)
		rrSetKey.PType = pool.RD
		_, rdPoolResponse, err := services.RecordService.Read(rrSetKey)

		if err == nil {
			if len(rdPoolResponse.RRSets) > 0 && rdPoolResponse.RRSets[0].OwnerName == rrSetKey.Owner {
				return errors.ResourceNotDestroyedError(rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccResourceRDPoolA(zoneName, ownerName string) string {
	return fmt.Sprintf(`
	%s
	resource "ultradns_rdpool" "a" {
		zone_name = "${resource.ultradns_zone.primary_rdpool.id}"
		owner_name = "%s"
		record_type = "1"
		ttl = 120
		record_data = ["192.168.1.1"]
		order = "FIXED"
		description = "RD Pool Resource of Type A"
	}
	`, acctest.TestAccResourceZonePrimary(zoneResourceName, zoneName), ownerName)
}

func testAccResourceUpdateRDPoolA(zoneName, ownerName string) string {
	return fmt.Sprintf(`
	%s
	resource "ultradns_rdpool" "a" {
		zone_name = "${resource.ultradns_zone.primary_rdpool.id}"
		owner_name = "%s.${resource.ultradns_zone.primary_rdpool.id}"
		record_type = "A"
		ttl = 150
		record_data = ["192.168.1.2"]
		order = "RANDOM"
	}
	`, acctest.TestAccResourceZonePrimary(zoneResourceName, zoneName), ownerName)
}

func testAccResourceRDPoolAAAA(zoneName, ownerName string) string {
	return fmt.Sprintf(`
	%s
	resource "ultradns_rdpool" "aaaa" {
		zone_name = "${resource.ultradns_zone.primary_rdpool.id}"
		owner_name = "%s"
		record_type = "AAAA"
		ttl = 120
		record_data = ["aaaa:bbbb:cccc:dddd:eeee:ffff:1111:2222"]
		order = "ROUND_ROBIN"
	}
	`, acctest.TestAccResourceZonePrimary(zoneResourceName, zoneName), ownerName)
}

func testAccResourceUpdateRDPoolAAAA(zoneName, ownerName string) string {
	return fmt.Sprintf(`
	%s
	resource "ultradns_rdpool" "aaaa" {
		zone_name = "${resource.ultradns_zone.primary_rdpool.id}"
		owner_name = "%s"
		record_type = "28"
		ttl = 150
		record_data = ["aaaa:bbbb:cccc:dddd:eeee:ffff:1111:3333"]
		order = "FIXED"
		description = "RD Pool Resource of Type AAAA"
	}
	`, acctest.TestAccResourceZonePrimary(zoneResourceName, zoneName), ownerName)
}
