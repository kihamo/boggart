package v3

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/mercury/v3"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{}
	ctx := r.Context()

	sn, buildDate, err := bind.provider.SerialNumberAndBuildDate()
	vars["serial_number"] = map[string]interface{}{
		"value": sn,
		"error": err,
	}
	vars["build_date"] = map[string]interface{}{
		"value": buildDate,
		"error": err,
	}

	fw, err := bind.provider.FirmwareVersion()
	vars["firmware_version"] = map[string]interface{}{
		"value": fw,
		"error": err,
	}

	address, err := bind.provider.Address()
	vars["address"] = map[string]interface{}{
		"value": address,
		"error": err,
	}

	valueFloat64, err := bind.provider.Frequency()
	vars["frequency"] = map[string]interface{}{
		"value": valueFloat64,
		"error": err,
	}

	phase1, phase2, phase3, err := bind.provider.Voltage()
	vars["voltage_phase1"] = map[string]interface{}{
		"value": phase1,
		"error": err,
	}
	vars["voltage_phase2"] = map[string]interface{}{
		"value": phase2,
		"error": err,
	}
	vars["voltage_phase3"] = map[string]interface{}{
		"value": phase3,
		"error": err,
	}

	phase1, phase2, phase3, err = bind.provider.Amperage()
	vars["amperage_phase1"] = map[string]interface{}{
		"value": phase1,
		"error": err,
	}
	vars["amperage_phase2"] = map[string]interface{}{
		"value": phase2,
		"error": err,
	}
	vars["amperage_phase3"] = map[string]interface{}{
		"value": phase3,
		"error": err,
	}

	valueFloat64, phase1, phase2, phase3, err = bind.provider.Power(v3.PowerNumberP)
	vars["power"] = map[string]interface{}{
		"value": valueFloat64,
		"error": err,
	}
	vars["power_phase1"] = map[string]interface{}{
		"value": phase1,
		"error": err,
	}
	vars["power_phase2"] = map[string]interface{}{
		"value": phase2,
		"error": err,
	}
	vars["power_phase3"] = map[string]interface{}{
		"value": phase3,
		"error": err,
	}

	phase1, phase2, phase3, err = bind.provider.PhasesAngle()
	vars["angle_phase1"] = map[string]interface{}{
		"value": phase1,
		"error": err,
	}
	vars["angle_phase2"] = map[string]interface{}{
		"value": phase2,
		"error": err,
	}
	vars["angle_phase3"] = map[string]interface{}{
		"value": phase3,
		"error": err,
	}

	valueFloat64, phase1, phase2, phase3, err = bind.provider.PowerFactors()
	vars["power_factors"] = map[string]interface{}{
		"value": valueFloat64,
		"error": err,
	}
	vars["power_factors_phase1"] = map[string]interface{}{
		"value": phase1,
		"error": err,
	}
	vars["power_factors_phase2"] = map[string]interface{}{
		"value": phase2,
		"error": err,
	}
	vars["power_factors_phase3"] = map[string]interface{}{
		"value": phase3,
		"error": err,
	}

	valueUint32, _, _, _, err := bind.provider.ReadArray(v3.ArrayReset, nil, v3.Tariff1)
	vars["power_t1"] = map[string]interface{}{
		"value": valueUint32,
		"error": err,
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
