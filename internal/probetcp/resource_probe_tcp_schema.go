package probetcp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/probe"
)

func resourceProbeTCPSchema() map[string]*schema.Schema {
	probeTCPSchema := probe.ResourceProbeSchema()

	probeTCPSchema["port"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  80,
	}

	probeTCPSchema["control_ip"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	probeTCPSchema["response"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem:     probe.SearchStringResource(),
	}

	probeTCPSchema["run_limit"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem:     probe.LimitResource(),
	}

	probeTCPSchema["avg_run_limit"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem:     probe.LimitResource(),
	}

	return probeTCPSchema
}
