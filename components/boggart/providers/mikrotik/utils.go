package mikrotik

func GetNameByMac(mac string, arp map[string]map[string]string, dns, dhcp map[string]string) (name string) {
	arpRecord, arpFound := arp[mac]
	if arpFound {
		if dnsName, ok := dns[arpRecord["address"]]; ok {
			name = dnsName
		}
	}

	if name == "" {
		if dhcpHost, ok := dhcp[mac]; ok {
			name = dhcpHost
		}
	}

	if name == "" && arpFound && arpRecord["comment"] != "" {
		return arpRecord["comment"]
	}

	return name
}
