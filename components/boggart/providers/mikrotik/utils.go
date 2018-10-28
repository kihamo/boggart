package mikrotik

func GetNameByMac(mac string, arp []IPARP, dns []IPDNSStatic, leases []IPDHCPServerLease) (name string) {
	var arpRow IPARP

	for _, row := range arp {
		if row.MacAddress == mac {
			for _, static := range dns {
				if static.Address == row.Address && !static.Disabled {
					name = static.Name
				}
			}

			arpRow = row
			break
		}
	}

	if name == "" {
		for _, lease := range leases {
			if lease.MacAddress == mac {
				name = lease.MacAddress
				break
			}
		}
	}

	if name == "" && arpRow.Comment != "" {
		return arpRow.Comment
	}

	return name
}
