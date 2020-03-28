package mikrotik

import (
	"github.com/kihamo/boggart/types"
)

func GetNameByMac(mac types.HardwareAddr, arp []IPARP, dns []IPDNSStatic, leases []IPDHCPServerLease) (name string) {
	var arpRow IPARP

	for _, row := range arp {
		if row.MacAddress.String() != "" && row.MacAddress.String() == mac.String() {
			for _, static := range dns {
				if static.Address.String() == row.Address.String() && !static.Disabled {
					name = static.Name
				}
			}

			arpRow = row

			break
		}
	}

	if name == "" {
		for _, lease := range leases {
			if lease.MacAddress.String() == mac.String() {
				name = lease.MacAddress.String()

				break
			}
		}
	}

	if name == "" && arpRow.Comment != "" {
		return arpRow.Comment
	}

	return name
}
