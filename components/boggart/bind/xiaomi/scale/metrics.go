package scale

import (
	"github.com/kihamo/snitch"
)

var (
	metricWeight            = snitch.NewGauge("weight_kilograms", "Weight in kilograms")
	metricImpedance         = snitch.NewGauge("impedance_ohm", "Impedance in ohm")
	metricBMR               = snitch.NewGauge("bmr", "BMR")
	metricBMI               = snitch.NewGauge("bmi", "BMI")
	metricFatPercentage     = snitch.NewGauge("fat_percentage", "Fat percentage")
	metricWaterPercentage   = snitch.NewGauge("water_percentage", "Water percentage")
	metricIdealWeight       = snitch.NewGauge("ideal_weight_kilograms", "Ideal weight in kilograms")
	metricLBMCoefficient    = snitch.NewGauge("lbm_coefficient", "LBM coefficient")
	metricBoneMass          = snitch.NewGauge("bone_mass_kilograms", "Bone mass in kilograms")
	metricMuscleMass        = snitch.NewGauge("muscle_mass_kilograms", "Muscle mass in kilograms")
	metricVisceralFat       = snitch.NewGauge("visceral_fat", "Visceral fat")
	metricFatMassToIdeal    = snitch.NewGauge("fat_mass_to_ideal", "Fat mass to ideal")
	metricProteinPercentage = snitch.NewGauge("protein_percentage", "Protein percentage")
	metricBodyType          = snitch.NewGauge("body_type", "Body type")
	metricMetabolicAge      = snitch.NewGauge("metabolic_age", "Metabolic age")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	mac := b.config.MAC.HardwareAddr.String()

	for name := range b.config.Profiles {
		metricWeight.With("mac", mac, "profile", name).Describe(ch)
		metricImpedance.With("mac", mac, "profile", name).Describe(ch)
		metricBMR.With("mac", mac, "profile", name).Describe(ch)
		metricBMI.With("mac", mac, "profile", name).Describe(ch)
		metricFatPercentage.With("mac", mac, "profile", name).Describe(ch)
		metricWaterPercentage.With("mac", mac, "profile", name).Describe(ch)
		metricIdealWeight.With("mac", mac, "profile", name).Describe(ch)
		metricLBMCoefficient.With("mac", mac, "profile", name).Describe(ch)
		metricBoneMass.With("mac", mac, "profile", name).Describe(ch)
		metricMuscleMass.With("mac", mac, "profile", name).Describe(ch)
		metricVisceralFat.With("mac", mac, "profile", name).Describe(ch)
		metricFatMassToIdeal.With("mac", mac, "profile", name).Describe(ch)
		metricProteinPercentage.With("mac", mac, "profile", name).Describe(ch)
		metricBodyType.With("mac", mac, "profile", name).Describe(ch)
		metricMetabolicAge.With("mac", mac, "profile", name).Describe(ch)
	}
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	mac := b.config.MAC.HardwareAddr.String()

	for name := range b.config.Profiles {
		metricWeight.With("mac", mac, "profile", name).Collect(ch)
		metricImpedance.With("mac", mac, "profile", name).Collect(ch)
		metricBMR.With("mac", mac, "profile", name).Collect(ch)
		metricBMI.With("mac", mac, "profile", name).Collect(ch)
		metricFatPercentage.With("mac", mac, "profile", name).Collect(ch)
		metricWaterPercentage.With("mac", mac, "profile", name).Collect(ch)
		metricIdealWeight.With("mac", mac, "profile", name).Collect(ch)
		metricLBMCoefficient.With("mac", mac, "profile", name).Collect(ch)
		metricBoneMass.With("mac", mac, "profile", name).Collect(ch)
		metricMuscleMass.With("mac", mac, "profile", name).Collect(ch)
		metricVisceralFat.With("mac", mac, "profile", name).Collect(ch)
		metricFatMassToIdeal.With("mac", mac, "profile", name).Collect(ch)
		metricProteinPercentage.With("mac", mac, "profile", name).Collect(ch)
		metricBodyType.With("mac", mac, "profile", name).Collect(ch)
		metricMetabolicAge.With("mac", mac, "profile", name).Collect(ch)
	}
}
