package mikrotik

func GetNameByMac(mac string, arp, dns, dhcp map[string]string) (name string) {
	if ip, ok := arp[mac]; ok {
		if dnsName, ok := dns[ip]; ok {
			name = dnsName
		}
	}

	if name == "" {
		if dhcpHost, ok := dhcp[mac]; ok {
			name = dhcpHost
		}
	}

	return name
}
