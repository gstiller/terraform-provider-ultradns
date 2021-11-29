package ultradns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/ultradns-go-sdk/ultradns"
)

func flattenPrimaryZone(zr *ultradns.ZoneResponse, rd *schema.ResourceData) *schema.Set {

	set := &schema.Set{
		F: schema.HashResource(primaryZoneCreateInfoResource()),
	}

	primaryCreateInfo := make(map[string]interface{})

	if val, ok := rd.GetOk("primary_create_info"); ok && val.(*schema.Set).Len() > 0 {
		primaryCreateInfo = val.(*schema.Set).List()[0].(map[string]interface{})
	}

	if len(zr.NotifyAddresses) > 0 {
		s := &schema.Set{
			F: schema.HashResource(notifyAddressResource()),
		}

		for _, notifyAddressData := range zr.NotifyAddresses {
			notifyAddress := make(map[string]interface{})

			notifyAddress["notify_address"] = notifyAddressData.NotifyAddress
			notifyAddress["description"] = notifyAddressData.Description

			s.Add(notifyAddress)
		}

		primaryCreateInfo["notify_addresses"] = s

	}

	if len(zr.RestrictIPList) > 0 {
		s := &schema.Set{
			F: schema.HashResource(restrictIpResource()),
		}

		for _, restrictIpData := range zr.RestrictIPList {
			restrictIp := make(map[string]interface{})

			restrictIp["start_ip"] = restrictIpData.StartIp
			restrictIp["end_ip"] = restrictIpData.EndIp
			restrictIp["cidr"] = restrictIpData.Cidr
			restrictIp["single_ip"] = restrictIpData.SingleIp
			restrictIp["comment"] = restrictIpData.Comment

			s.Add(restrictIp)
		}

		primaryCreateInfo["restrict_ip"] = s

	}

	if zr.Tsig != nil {
		tsig := make(map[string]interface{})

		tsig["tsig_key_name"] = zr.Tsig.TsigKeyName
		tsig["tsig_key_value"] = zr.Tsig.TsigKeyValue
		tsig["tsig_algorithm"] = zr.Tsig.TsigAlgorithm
		tsig["description"] = zr.Tsig.Description

		s := &schema.Set{
			F: schema.HashResource(tsigResource()),
		}

		s.Add(tsig)

		primaryCreateInfo["tsig"] = s
	}

	set.Add(primaryCreateInfo)

	return set
}

func flattenSecondaryZone(zr *ultradns.ZoneResponse, rd *schema.ResourceData) *schema.Set {
	set := &schema.Set{
		F: schema.HashResource(secondaryZoneCreateInfoResource()),
	}

	secondaryCreateInfo := make(map[string]interface{})

	if val, ok := rd.GetOk("secondary_create_info"); ok && val.(*schema.Set).Len() > 0 {
		secondaryCreateInfo = val.(*schema.Set).List()[0].(map[string]interface{})
	}

	if zr.NotificationEmailAddress != "" {
		secondaryCreateInfo["notification_email_address"] = zr.NotificationEmailAddress
	}

	if zr.PrimaryNameServers != nil {
		s := &schema.Set{
			F: schema.HashResource(nameServerResource()),
		}
		if zr.PrimaryNameServers.NameServerIpList != nil {
			if zr.PrimaryNameServers.NameServerIpList.NameServerIp1 != nil {
				s.Add(getNameServer(zr.PrimaryNameServers.NameServerIpList.NameServerIp1))
			}
			if zr.PrimaryNameServers.NameServerIpList.NameServerIp2 != nil {
				s.Add(getNameServer(zr.PrimaryNameServers.NameServerIpList.NameServerIp2))
			}
			if zr.PrimaryNameServers.NameServerIpList.NameServerIp3 != nil {
				s.Add(getNameServer(zr.PrimaryNameServers.NameServerIpList.NameServerIp3))
			}
		}
		secondaryCreateInfo["primary_name_server"] = s
	}
	set.Add(secondaryCreateInfo)

	return set
}

func flattenAliasZone(zr *ultradns.ZoneResponse, rd *schema.ResourceData) *schema.Set {
	set := &schema.Set{
		F: schema.HashResource(aliasZoneCreateInfoResource()),
	}

	aliasCreateInfo := make(map[string]interface{})
	aliasCreateInfo["original_zone_name"] = zr.OriginalZoneName

	set.Add(aliasCreateInfo)

	return set
}

func getNameServer(ns *ultradns.NameServerIp) map[string]interface{} {
	nameserver := make(map[string]interface{})
	nameserver["ip"] = ns.Ip
	nameserver["tsig_key"] = ns.TsigKey
	nameserver["tsig_key_value"] = ns.TsigKeyValue
	nameserver["tsig_algorithm"] = ns.TsigAlgorithm
	return nameserver
}

func validateZoneName(i interface{}, s string) (warns []string, errs []error) {
	zoneName := i.(string)
	if len(zoneName) > 0 {
		if lastChar := zoneName[len(zoneName)-1]; lastChar != '.' {
			errs = append(errs, fmt.Errorf("zone name must be a FQDN"))
		}
	}
	return
}
