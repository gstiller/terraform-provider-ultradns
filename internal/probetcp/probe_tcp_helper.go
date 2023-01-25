package probetcp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/probe"
	sdkprobe "github.com/ultradns/ultradns-go-sdk/pkg/probe"
	"github.com/ultradns/ultradns-go-sdk/pkg/probe/tcp"
)

func flattenProbeTCP(probeData *sdkprobe.Probe, rd *schema.ResourceData) error {
	if err := probe.FlattenProbe(probeData, rd); err != nil {
		return err
	}

	details := probeData.Details.(*tcp.Details)

	if err := rd.Set("port", details.Port); err != nil {
		return err
	}

	if err := rd.Set("control_ip", details.ControlIP); err != nil {
		return err
	}

	if details.Limits != nil {
		limits := details.Limits

		if err := rd.Set("response", probe.GetSearchStringList(limits.Response)); err != nil {
			return err
		}

		if err := rd.Set("run_limit", probe.GetLimitList(limits.Run)); err != nil {
			return err
		}

		if err := rd.Set("avg_run_limit", probe.GetLimitList(limits.AvgRun)); err != nil {
			return err
		}
	}

	return nil
}
