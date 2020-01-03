package scale

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricWeight            = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_weight_kilograms", "Weight in kilograms")
	metricImpedance         = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_impedance_ohm", "Impedance in ohm")
	metricBMR               = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_bmr", "BMR")
	metricBMI               = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_bmi", "BMI")
	metricFatPercentage     = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_fat_percentage", "Fat percentage")
	metricWaterPercentage   = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_water_percentage", "Water percentage")
	metricIdealWeight       = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_ideal_weight_kilograms", "Ideal weight in kilograms")
	metricLBMCoefficient    = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_lbm_coefficient", "LBM coefficient")
	metricBoneMass          = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_bone_mass_kilograms", "Bone mass in kilograms")
	metricMuscleMass        = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_muscle_mass_kilograms", "Muscle mass in kilograms")
	metricVisceralFat       = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_visceral_fat", "Visceral fat")
	metricFatMassToIdeal    = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_fat_mass_to_ideal", "Fat mass to ideal")
	metricProteinPercentage = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_protein_percentage", "Protein percentage")
	metricBodyType          = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_body_type", "Body type")
	metricMetabolicAge      = snitch.NewGauge(boggart.ComponentName+"_bind_xiaomi_scale_metabolic_age", "Metabolic age")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	for name := range b.config.Profiles {
		metricWeight.With("serial_number", sn, "profile", name).Describe(ch)
		metricImpedance.With("serial_number", sn, "profile", name).Describe(ch)
		metricBMR.With("serial_number", sn, "profile", name).Describe(ch)
		metricBMI.With("serial_number", sn, "profile", name).Describe(ch)
		metricFatPercentage.With("serial_number", sn, "profile", name).Describe(ch)
		metricWaterPercentage.With("serial_number", sn, "profile", name).Describe(ch)
		metricIdealWeight.With("serial_number", sn, "profile", name).Describe(ch)
		metricLBMCoefficient.With("serial_number", sn, "profile", name).Describe(ch)
		metricBoneMass.With("serial_number", sn, "profile", name).Describe(ch)
		metricMuscleMass.With("serial_number", sn, "profile", name).Describe(ch)
		metricVisceralFat.With("serial_number", sn, "profile", name).Describe(ch)
		metricFatMassToIdeal.With("serial_number", sn, "profile", name).Describe(ch)
		metricProteinPercentage.With("serial_number", sn, "profile", name).Describe(ch)
		metricBodyType.With("serial_number", sn, "profile", name).Describe(ch)
		metricMetabolicAge.With("serial_number", sn, "profile", name).Describe(ch)
	}
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	for name := range b.config.Profiles {
		metricWeight.With("serial_number", sn, "profile", name).Collect(ch)
		metricImpedance.With("serial_number", sn, "profile", name).Collect(ch)
		metricBMR.With("serial_number", sn, "profile", name).Collect(ch)
		metricBMI.With("serial_number", sn, "profile", name).Collect(ch)
		metricFatPercentage.With("serial_number", sn, "profile", name).Collect(ch)
		metricWaterPercentage.With("serial_number", sn, "profile", name).Collect(ch)
		metricIdealWeight.With("serial_number", sn, "profile", name).Collect(ch)
		metricLBMCoefficient.With("serial_number", sn, "profile", name).Collect(ch)
		metricBoneMass.With("serial_number", sn, "profile", name).Collect(ch)
		metricMuscleMass.With("serial_number", sn, "profile", name).Collect(ch)
		metricVisceralFat.With("serial_number", sn, "profile", name).Collect(ch)
		metricFatMassToIdeal.With("serial_number", sn, "profile", name).Collect(ch)
		metricProteinPercentage.With("serial_number", sn, "profile", name).Collect(ch)
		metricBodyType.With("serial_number", sn, "profile", name).Collect(ch)
		metricMetabolicAge.With("serial_number", sn, "profile", name).Collect(ch)
	}
}
