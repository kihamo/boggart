package z_stack

type definition struct {
	name  string
	types []uint16
}

var definitions = map[uint16]map[uint16]definition{
	SubSystemReserved: {},
	SubSystemSysInterface: {
		0x00: {
			name:  "SYS_RESET_REQ",
			types: []uint16{TypeAREQ},
		},
		0x01: {
			name:  "SYS_PING",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x02: {
			name:  "SYS_VERSION",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x03: {
			name:  "SYS_SET_EXTADDR",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x04: {
			name:  "SYS_GET_EXTADDR",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x05: {
			name:  "SYS_RAM_READ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x06: {
			name:  "SYS_RAM_WRITE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x07: {
			name:  "SYS_OSAL_NV_ITEM_INIT",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x08: {
			name:  "SYS_OSAL_NV_READ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x09: {
			name:  "SYS_OSAL_NV_WRITE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0A: {
			name:  "SYS_OSAL_START_TIMER",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0B: {
			name:  "SYS_OSAL_STOP_TIMER",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0C: {
			name:  "SYS_RANDOM",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0D: {
			name:  "SYS_ADC_READ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0E: {
			name:  "SYS_GPIO",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0F: {
			name:  "SYS_STACK_TUNE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x80: {
			name:  "SYS_RESET_IND",
			types: []uint16{TypeAREQ},
		},
	},
	SubSystemMACInterface: {
		0x01: {
			name:  "MT_MAC_RESET_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x02: {
			name:  "MT_MAC_INIT",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x03: {
			name:  "MT_MAC_START_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x04: {
			name:  "MT_MAC_SYNC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x05: {
			name:  "MT_MAC_DATA_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x06: {
			name:  "MT_MAC_ASSOCIATE_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x50: {
			name:  "MT_MAC_ASSOCIATE_RSP",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x07: {
			name:  "MT_MAC_DISASSOCIATE_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x08: {
			name:  "MT_MAC_GET_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x09: {
			name:  "MT_MAC_SET_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0C: {
			name:  "MT_MAC_SCAN_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x51: {
			name:  "MT_MAC_ORPHAN_RSP",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0D: {
			name:  "MT_MAC_POLL_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0E: {
			name:  "MT_MAC_PURGE_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0F: {
			name:  "MT_MAC_SET_RX_GAIN_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x10: {
			name:  "MT_MAC_SRC_MATCH_ENABLE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x11: {
			name:  "MT_MAC_SRC_MATCH_ADD_ENTRY",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x12: {
			name:  "MT_MAC_SRC_MATCH_DEL_ENTRY",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x13: {
			name:  "MT_MAC_SRC_MATCH_CHECK_SRC_ADDR",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x14: {
			name:  "MT_MAC_SRC_MATCH_ACK_ALL_PENDING",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x15: {
			name:  "MT_MAC_SRC_MATCH_CHECK_ALL_PENDING",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x80: {
			name:  "MT_MAC_SYNC_LOSS_IND",
			types: []uint16{TypeAREQ},
		},
		0x81: {
			name:  "MT_MAC_ASSOCIATE_IND",
			types: []uint16{TypeAREQ},
		},
		0x82: {
			name:  "MT_MAC_ASSOCIATE_CNF",
			types: []uint16{TypeAREQ},
		},
		0x83: {
			name:  "MT_MAC_BEACON_NOTIFY_IND",
			types: []uint16{TypeAREQ},
		},
		0x84: {
			name:  "MT_MAC_DATA_CNF",
			types: []uint16{TypeAREQ},
		},
		0x85: {
			name:  "MT_MAC_DATA_IND",
			types: []uint16{TypeAREQ},
		},
		0x86: {
			name:  "MT_MAC_DISASSOCIATE_IND",
			types: []uint16{TypeAREQ},
		},
		0x87: {
			name:  "MT_MAC_DISASSOCIATE_CNF",
			types: []uint16{TypeAREQ},
		},
		0x8A: {
			name:  "MT_MAC_ORPHAN_IND",
			types: []uint16{TypeAREQ},
		},
		0x8B: {
			name:  "MT_MAC_POLL_CNF",
			types: []uint16{TypeAREQ},
		},
		0x8C: {
			name:  "MT_MAC_SCAN_CNF",
			types: []uint16{TypeAREQ},
		},
		0x8D: {
			name:  "MT_MAC_COMM_STATUS_IND",
			types: []uint16{TypeAREQ},
		},
		0x8E: {
			name:  "MT_MAC_START_CNF",
			types: []uint16{TypeAREQ},
		},
		0x8F: {
			name:  "MT_MAC_RX_ENABLE_CNF",
			types: []uint16{TypeAREQ},
		},
		0x9A: {
			name:  "MT_MAC_PURGE_CNF",
			types: []uint16{TypeAREQ},
		},
	},
	SubSystemNWKInterface: {},
	SubSystemAFInterface: {
		0x00: {
			name:  "AF_REGISTER",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x80: {
			name:  "AF_DATA_CONFIRM",
			types: []uint16{TypeAREQ},
		},
		0x81: {
			name:  "AF_INCOMING_MSG",
			types: []uint16{TypeAREQ},
		},
		0x82: {
			name:  "AF_INCOMING_MSG_EXT",
			types: []uint16{TypeAREQ},
		},
	},
	SubSystemZDOInterface: {
		0x00: {
			name:  "ZDO_NWK_ADDR_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x01: {
			name:  "ZDO_IEEE_ADDR_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x02: {
			name:  "ZDO_NODE_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x03: {
			name:  "ZDO_POWER_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x04: {
			name:  "ZDO_SIMPLE_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x05: {
			name:  "ZDO_ACTIVE_EP_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x06: {
			name:  "ZDO_MATCH_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x07: {
			name:  "ZDO_COMPLEX_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x08: {
			name:  "ZDO_USER_DESC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0A: {
			name:  "ZDO_END_DEVICE_ANNCE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0B: {
			name:  "ZDO_USER_DESC_SET",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0C: {
			name:  "ZDO_SERVER_DISC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x20: {
			name:  "ZDO_END_DEVICE_BIND_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x21: {
			name:  "ZDO_BIND_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x22: {
			name:  "ZDO_UNBIND_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x30: {
			name:  "ZDO_MGMT_NWK_DISC_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x31: {
			name:  "ZDO_MGMT_LQI_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x32: {
			name:  "ZDO_MGMT_RTG_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x33: {
			name:  "ZDO_MGMT_BIND_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x34: {
			name:  "ZDO_MGMT_LEAVE_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x35: {
			name:  "ZDO_MGMT_DIRECT_JOIN_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x36: {
			name:  "ZDO_MGMT_PERMIT_JOIN_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x37: {
			name:  "ZDO_MGMT_NWK_UPDATE_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x3E: {
			name:  "ZDO_MSG_CB_REGISTER",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x3F: {
			name:  "ZDO_MSG_CB_REMOVE",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x40: {
			name:  "ZDO_STARTUP_FROM_APP",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x41: {
			name:  "ZDO_AUTO_FIND_DESTINATION",
			types: []uint16{TypeAREQ},
		},
		0x4A: {
			name:  "UTIL_ASSOC_GET_WITH_ADDRESS",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x23: {
			name:  "ZDO_SET_LINK_KEY",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x24: {
			name:  "ZDO_REMOVE_LINK_KEY",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x25: {
			name:  "ZDO_GET_LINK_KEY",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x26: {
			name:  "ZDO_NETWORK_DISCOVERY_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x27: {
			name:  "ZDO_JOIN_REQ",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x80: {
			name:  "ZDO_NWK_ADDR_RSP",
			types: []uint16{TypeAREQ},
		},
		0x81: {
			name:  "ZDO_IEEE_ADDR_RSP",
			types: []uint16{TypeAREQ},
		},
		0x82: {
			name:  "ZDO_NODE_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x83: {
			name:  "ZDO_POWER_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x84: {
			name:  "ZDO_SIMPLE_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x85: {
			name:  "ZDO_ACTIVE_EP_RSP",
			types: []uint16{TypeAREQ},
		},
		0x86: {
			name:  "ZDO_MATCH_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x87: {
			name:  "ZDO_COMPLEX_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x88: {
			name:  "ZDO_USER_DESC_RSP",
			types: []uint16{TypeAREQ},
		},
		0x89: {
			name:  "ZDO_USER_DESC_CONF",
			types: []uint16{TypeAREQ},
		},
		0x8A: {
			name:  "ZDO_SERVER_DISC_RSP",
			types: []uint16{TypeAREQ},
		},
		0xA0: {
			name:  "ZDO_END_DEVICE_BIND_RSP",
			types: []uint16{TypeAREQ},
		},
		0xA1: {
			name:  "ZDO_BIND_RSP",
			types: []uint16{TypeAREQ},
		},
		0xA2: {
			name:  "ZDO_UNBIND_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB0: {
			name:  "ZDO_MGMT_NWK_DISC_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB1: {
			name:  "ZDO_MGMT_LQI_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB2: {
			name:  "ZDO_MGMT_RTG_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB3: {
			name:  "ZDO_MGMT_BIND_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB4: {
			name:  "ZDO_MGMT_LEAVE_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB5: {
			name:  "ZDO_MGMT_DIRECT_JOIN_RSP",
			types: []uint16{TypeAREQ},
		},
		0xB6: {
			name:  "ZDO_MGMT_PERMIT_JOIN_RSP",
			types: []uint16{TypeAREQ},
		},
		0xC0: {
			name:  "ZDO_STATE_CHANGE_IND",
			types: []uint16{TypeAREQ},
		},
		0xC1: {
			name:  "ZDO_END_DEVICE_ANNCE_IND",
			types: []uint16{TypeAREQ},
		},
		0xC2: {
			name:  "ZDO_MATCH_DESC_RSP_SENT",
			types: []uint16{TypeAREQ},
		},
		0xC3: {
			name:  "ZDO_STATUS_ERROR_RSP",
			types: []uint16{TypeAREQ},
		},
		0xC4: {
			name:  "ZDO_SRC_RTG_IND",
			types: []uint16{TypeAREQ},
		},
		0xC5: {
			name:  "ZDO_BEACON_NOTIFY_IND",
			types: []uint16{TypeAREQ},
		},
		0xC6: {
			name:  "ZDO_JOIN_CNF",
			types: []uint16{TypeAREQ},
		},
		0xC7: {
			name:  "ZDO_NWK_DISCOVERY_CNF",
			types: []uint16{TypeAREQ},
		},
		0xC8: {
			name:  "ZDO_CONCENTRATOR_IND_CB",
			types: []uint16{TypeAREQ},
		},
		0xFF: {
			name:  "ZDO_MSG_CB_INCOMING",
			types: []uint16{TypeAREQ},
		},
	},
	SubSystemSAPIInterface: {
		0x05: {
			name:  "ZB_WRITE_CONFIGURATION",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
	},
	SubSystemUtilInterface: {
		0x00: {
			name:  "UTIL_GET_DEVICE_INFO",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
		0x0A: {
			name:  "UTIL_LED_CONTROL",
			types: []uint16{TypeSREQ, TypeSRSP},
		},
	},
	SubSystemDebugInterface: {},
	SubSystemAppInterface:   {},
}
