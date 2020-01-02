package scale

// https://github.com/PotatoSpudowski/L.I.S.A/blob/268edfbbec1ee84a68d1cf2d8de7e30720b66cac/Body_Metrics.py
// https://github.com/wiecosystem/Bluetooth/blob/master/sandbox/huami.health.scale2/body_metrics.py

type sex uint8

type Unit byte

const (
	SexMale sex = iota
	SexFemale

	UnitKG  Unit = 0x02
	UnitKG2 Unit = 0x22
	UnitLBS Unit = 0x03
	UnitJIN Unit = 0x12
)

func (u Unit) String() string {
	switch u {
	case UnitKG, UnitKG2:
		return "kg"

	case UnitLBS:
		return "lbs"

	case UnitJIN:
		return "jin"

	default:
		return "unknown"
	}
}

func checkValueOverflow(value, min, max float64) float64 {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}
