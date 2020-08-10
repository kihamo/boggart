package z_stack

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"sync/atomic"
	"time"
)

var endpoints = []Endpoint{
	{EndPoint: 1, AppProfId: 0x0104},
	{EndPoint: 2, AppProfId: 0x0101},
	{EndPoint: 3, AppProfId: 0x0105},
	{EndPoint: 4, AppProfId: 0x0107},
	{EndPoint: 5, AppProfId: 0x0108},
	{EndPoint: 6, AppProfId: 0x0109},
	{EndPoint: 8, AppProfId: 0x0104},
	{EndPoint: 11, AppProfId: 0x0104, AppDeviceId: 0x0400, AppOutClusterList: []uint16{1280, 1282}},
	{EndPoint: 12, AppProfId: 0xC05E},
	{EndPoint: 13, AppProfId: 0x0104, AppInClusterList: []uint16{25}},
	{EndPoint: 47, AppProfId: 0x0104},
	{EndPoint: 110, AppProfId: 0x0104},
	{EndPoint: 242, AppProfId: 0xA1E0},
}

type NvItem struct {
	id    uint16
	value []byte
}

var (
	nvItemZnpHasConfigured = func(version *SysVersion) NvItem {
		i := NvItem{
			id:    NvItemIDHasConfiguredZStack1,
			value: []byte{0x55},
		}

		if version.Product != VersionZStack12 {
			i.id = NvItemIDHasConfiguredZStack3
		}

		return i
	}
	nvItemStartupOption = func() NvItem {
		return NvItem{
			id:    0x03,
			value: []byte{0x02},
		}
	}
	nvItemLogicalType = func(t uint8) NvItem {
		return NvItem{
			id:    0x87,
			value: []byte{t},
		}
	}
	nvItemNetworkKeyDistribute = func(distribute bool) NvItem {
		i := NvItem{
			id: 0x63,
		}

		if distribute {
			i.value = []byte{0x01}
		} else {
			i.value = []byte{0x00}
		}

		return i
	}
	nvItemZdoDirectCb = func() NvItem {
		return NvItem{
			id:    0x8F,
			value: []byte{0x01},
		}
	}
	nvItemPanID = func(panID uint16) NvItem {
		value := make([]byte, 2)
		binary.LittleEndian.PutUint16(value, panID)

		return NvItem{
			id:    0x83,
			value: value,
		}
	}
	nvItemChannelList = func(channel uint32) NvItem {
		channel = 1 << channel

		value := make([]byte, 4)
		binary.LittleEndian.PutUint32(value, channel)

		return NvItem{
			id:    0x84,
			value: value,
		}
	}
	nvItemExtendedPanID = func(panID []byte) NvItem {
		return NvItem{
			id:    0x2D,
			value: panID,
		}
	}
	nvItemTcLinKey = func() NvItem {
		return NvItem{
			id: 0x0101,
			value: []byte{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x5A, 0x69, 0x67, 0x42, 0x65, 0x65, 0x41, 0x6C,
				0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		}
	}
)

func (c *Client) Version(ctx context.Context) (_ *SysVersion, err error) {
	c.versionOnce.Do(func() {
		c.version, err = c.SysVersion(ctx)
	})

	if err != nil {
		c.versionOnce.Reset()
	}

	return c.version, err
}

func (c *Client) SkipBootLoader() error {
	_, err := c.Write([]byte{0xEF})
	return err
}

func (c *Client) Reset(ctx context.Context, t uint8) error {
	/*
		SYS_RESET_IND

		This command is sent by the device to indicate the reset

		Usage:
			AREQ:
				       1      |      1      |      1      |    1   |       1      |     1     |     1    |    1     |   1
				Length = 0x06 | Cmd0 = 0x41 | Cmd1 = 0x80 | Reason | TransportRev | ProductId | MajorRel | MinorRel | HwRev
			Attributes:
				Reason       1 byte Reason for the reset.
				                    0x00 Power-up
				                    0x01 External
				                    0x02 Watch-dog
				TransportRev 1 byte Transport protocol revision.
				Product      1 byte Major release number.
				MinorRel     1 byte Minor release number.
				HwRev        1 byte Hardware revision number.
	*/
	waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
		return response.Type() == TypeAREQ && response.SubSystem() == SubSystemSysInterface && response.CommandID() == 0x80
	})

	err := c.SysReset(ctx, t)
	if err != nil {
		return err
	}

	select {
	case <-waitResponse:
		// TODO: return struct
		return nil

	case e := <-waitErr:
		return e
	}

	return nil
}

func (c *Client) Boot(ctx context.Context) (err error) {
	if err := c.SkipBootLoader(); err != nil {
		return err
	}
	time.Sleep(time.Second)

	// Initialization check
	valid := true

	// check hasConfigured flag
	if valid, err = c.InitializationCheck(ctx); err != nil {
		return err
	}

	// Initialization run if valid == false
	if !valid {
		if err = c.initialization(ctx); err != nil {
			return err
		}
	}

	/*
	 * Configuration
	 */

	/*
	 * Start as coordinator
	 */
	device, err := c.UtilGetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	if device.DeviceState != DeviceStateStartedCoordinator {
		waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
			return response.Type() == TypeAREQ &&
				response.SubSystem() == SubSystemZDOInterface &&
				response.CommandID() == 0xC0 &&
				DeviceState(response.Data()[0]) == DeviceStateStartedCoordinator
		})

		status, err := c.ZDOStartupFromApp(ctx, 100)
		if err != nil {
			return err
		}

		if status != 0 && status != 1 {
			return errors.New("startup from app failed")
		}

		select {
		case <-waitResponse:

		case e := <-waitErr:
			return e
		}
	}

	// регистрируем привязанные устройства хотя инормация по ним не полная
	// полную информацию добираем через синхронизацию с таблицей на устройстве
	if len(device.AssocDevicesList) > 0 {
		for _, networkAddress := range device.AssocDevicesList {
			d := NewDevice()
			d.SetNetworkAddress(networkAddress)

			c.deviceAdd(d)
		}

		go func() {
			c.SyncDevices(context.Background())
		}()
	}

	/*
	 * Register endpoints
	 */
	/*
		ZDO_ACTIVE_EP_RSP

		This callback message is in response to the ZDO Active Endpoint Request.

		Usage:
			AREQ:
				         1         |       1     |       1     |    2    |   1    |    2    |       1       |     0-77
				Length = 0x06-0x53 | Cmd0 = 0x45 | Cmd1 = 0x85 | SrcAddr | Status | NwkAddr | ActiveEPCount | ActiveEPList
			Attributes:
				SrcAddr       2 bytes    The message’s source network address.
				Status        1 bytes    This field indicates either SUCCESS or FAILURE.
				NWKAddr       2 bytes    Device’s short address that this response describes.
				ActiveEPCount 1 byte     Number of active endpoint in the list
				ActiveEPList  0-77 bytes Array of active endpoints on this device.

		Example from zigbee2mqtt:
			zigbee-herdsman:adapter:zStack:znp:AREQ <-- ZDO - activeEpRsp - {"srcaddr":0,"status":0,"nwkaddr":0,"activeepcount":0,"activeeplist":[]} +14ms
	*/
	waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
		return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == CommandActiveEndpointResponse
	})

	err = c.ZDOActiveEndpoints(ctx)
	if err != nil {
		return err
	}

	registeredEndpoints := make(map[uint8]bool)

	select {
	case r := <-waitResponse:
		dataOut := r.Data()

		if dataOut[2] != 0 {
			return errors.New("get active endpoints failed")
		}

		for _, id := range dataOut[6:] {
			registeredEndpoints[uint8(id)] = true
		}

	case e := <-waitErr:
		return e
	}

	// register endpoints
	for _, endpoint := range endpoints {
		if !registeredEndpoints[endpoint.EndPoint] {
			if err := c.AfRegister(ctx, endpoint); err != nil {
				return err
			}
		}
	}

	/*
	 * Group green power
	 */
	ep := uint8(242)
	groupID := uint16(0x0B84)

	// check exists
	if group, err := c.ZDOExtFindGroup(ctx, ep, groupID); err != nil {
		return err
	} else if group.Status != CommandStatusSuccess {
		// register if not exists
		if err = c.ZDOExtAddToGroup(ctx, ep, groupID, nil); err != nil {
			return err
		}
	}

	/*
	 * Start default watcher
	 */

	go c.defaultWatcher()

	return nil
}

func (c *Client) PermitJoinEnabled() bool {
	return atomic.LoadUint32(&c.permitJoin) != 0
}

func (c *Client) PermitJoinDisable(ctx context.Context) error {
	return c.PermitJoin(ctx, 0)
}

func (c *Client) PermitJoin(ctx context.Context, seconds uint8) error {
	if (seconds > 0 && c.PermitJoinEnabled()) || (seconds == 0 && !c.PermitJoinEnabled()) {
		return nil
	}

	/*
		ZDO_MGMT_PERMIT_JOIN_RSP

		This callback message is in response to the ZDO Management Permit Join Request.

		Usage:
			AREQ:
				       1      |       1     |       1     |    2    |   1
				Length = 0x03 | Cmd0 = 0x45 | Cmd1 = 0xB6 | SrcAddr | Status
			Attributes:
				SrcAddr       2 bytes    Source address of the message.
				Status        1 bytes    This field indicates either SUCCESS (0) or FAILURE (1).

		Example from zigbee2mqtt:
			zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,3,69,182,0,0,0,240] +9ms
			zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,3,69,182,0,0,0,240] +1ms
			zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 3 - 2 - 5 - 182 - [0,0,0] - 240 +0ms
			zigbee-herdsman:adapter:zStack:znp:AREQ <-- ZDO - mgmtPermitJoinRsp - {"srcaddr":0,"status":0} +47ms
	*/

	waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
		return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == CommandManagementPermitJoinResponse
	})

	/*
		Судя по коду zigbee2mqtt в 3 версии протокола 255 (постоянно включено) установить нельзя, так
		происходит защита сети, поэтому для 3 версии протокола включается специальный механизм который
		устанавливает 254 секундный интервал и переактивирует его по истечению этого времени. В версии
		1.2 протокола можно установить 255 то есть активировать постоянно. Поэтому тут TODO сделать хак
		но пока не актуально так как стик 1.2

		В любых версия по истечению интервала (а при 255 интервале сразу после установки) приходит
		пакет с CmdID=0xCB (permitJoinInd), который оповещает что время истекло. В теле пакета
		{name: 'duration', parameterType: ParameterType.UINT8} который содержит интервал который был
		установлен ранее.
	*/

	if err := c.ZDOPermitJoin(ctx, seconds); err != nil {
		return err
	}

	select {
	case r := <-waitResponse:
		if r.Data()[2] != 0 {
			return errors.New("enable permit join failed")
		}

		if seconds == 0 {
			atomic.StoreUint32(&c.permitJoin, 0)
		} else {
			atomic.StoreUint32(&c.permitJoin, 1)
		}

	case e := <-waitErr:
		return e
	}

	return nil
}

func (c *Client) SyncDevices(ctx context.Context) error {
	// вычитывает таблицу LQI с устройства, эта операция медленная поэтому запускается в фоне
	// чтобы собрать детальную информацию про устройства

	list, err := c.LQI(ctx, 0)
	if err != nil {
		return err
	}

	for _, item := range list {
		if device := c.Device(item.NetworkAddress); device != nil {
			device.SetIEEEAddress(item.ExtendedAddress)
			device.SetDeviceType(item.DeviceType)
		}
	}

	return nil
}

func (c *Client) LQI(ctx context.Context, networkAddress uint16) ([]NeighborLqiListItem, error) {
	request := func(index uint8) (*ZDOLQIMessage, error) {
		waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
			return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == 0xB1
		})

		if err := c.ZDOLQI(ctx, networkAddress, index); err != nil {
			return nil, err
		}

		for {
			select {
			case frame := <-waitResponse:
				return c.ZDOLQIMessage(frame)

			case e := <-waitErr:
				return nil, e
			}
		}
	}

	msg, err := request(0)
	if err != nil {
		return nil, err
	}

	total := int(msg.NeighborTableEntries)
	list := make([]NeighborLqiListItem, 0, total)
	list = append(list, msg.NeighborLqiList...)

	for i := uint8(1); len(list) < total; i++ {
		msg, err := request(i)
		if err != nil {
			return nil, err
		}

		list = append(list, msg.NeighborLqiList...)
	}

	return list, nil
}

func (c *Client) NetworkDiscovery(ctx context.Context) ([]NetworkListItem, error) {
	request := func(index uint8) (*ZDOManagementNetworkDiscoveryMessage, error) {
		scanDuration := uint8(1)

		waitResponse, waitErr := c.WaitAsyncWithTimeout(ctx, func(response *Frame) bool {
			return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == 0xB0
		}, time.Duration(scanDuration+6)*time.Second)

		if err := c.ZDOManagementNetworkDiscovery(ctx, 0, ScanChannelsAllChannels, scanDuration, index); err != nil {
			return nil, err
		}

		for {
			select {
			case frame := <-waitResponse:
				return c.ZDONetworkDiscoveryMessage(frame)

			case e := <-waitErr:
				return nil, e
			}
		}
	}

	msg, err := request(0)
	if err != nil {
		return nil, err
	}

	total := int(msg.NetworkCount)
	list := make([]NetworkListItem, 0, total)
	list = append(list, msg.NetworkList...)

	for i := uint8(1); len(list) < total; i++ {
		msg, err := request(i)
		if err != nil {
			return nil, err
		}

		list = append(list, msg.NetworkList...)
	}

	return list, nil
}

func (c *Client) RoutingTable(ctx context.Context, dstAddr uint16) ([]RoutingTableListItem, error) {
	request := func(index uint8) (*ZDOManagementRoutingTableMessage, error) {
		waitResponse, waitErr := c.WaitAsync(ctx, func(response *Frame) bool {
			return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == CommandManagementRoutingTableResponse
		})

		if err := c.ZDORoutingTable(ctx, dstAddr, index); err != nil {
			return nil, err
		}

		for {
			select {
			case frame := <-waitResponse:
				return c.ZDOManagementRoutingTableMessage(frame)

			case e := <-waitErr:
				return nil, e
			}
		}
	}

	msg, err := request(0)
	if err != nil {
		return nil, err
	}

	total := int(msg.RoutingTableListCount)
	list := make([]RoutingTableListItem, 0, total)
	list = append(list, msg.RoutingTableList...)

	for i := uint8(1); len(list) < total; i++ {
		msg, err := request(i)
		if err != nil {
			return nil, err
		}

		list = append(list, msg.RoutingTableList...)
	}

	return list, nil
}

func (c *Client) LEDEnabled() bool {
	return atomic.LoadUint32(&c.enabledLed) != 0
}

func (c *Client) LEDSupport(ctx context.Context) bool {
	version, err := c.Version(ctx)
	if err != nil {
		return false
	}

	return version.Product != VersionZStack3x0
}

func (c *Client) LEDDisable(ctx context.Context) error {
	return c.LED(ctx, false)
}

func (c *Client) LED(ctx context.Context, enabled bool) error {
	if !c.LEDSupport(ctx) {
		return errors.New("adapter doesn't support LED")
	}

	err := c.UtilLEDControl(ctx, 3, enabled)
	if err == nil {
		if enabled {
			atomic.StoreUint32(&c.enabledLed, 1)
		} else {
			atomic.StoreUint32(&c.enabledLed, 0)
		}
	}

	return err
}

func (c *Client) nwItemValidate(ctx context.Context, item NvItem) (bool, error) {
	response, err := c.SysOsalNvRead(ctx, item.id, 0x00)
	if err != nil {
		return false, err
	}

	return bytes.Equal(response.Value, item.value), nil
}

func (c *Client) nvItemWrite(ctx context.Context, item NvItem) error {
	return c.SysOsalNvWrite(ctx, item.id, 0, item.value)
}

func (c *Client) InitializationCheck(ctx context.Context) (valid bool, err error) {
	version, err := c.Version(ctx)
	if err != nil {
		return false, err
	}

	valid = true

	if valid, err = c.nwItemValidate(ctx, nvItemZnpHasConfigured(version)); err != nil {
		return false, err
	}

	if valid {
		if valid, err = c.nwItemValidate(ctx, nvItemChannelList(c.channel)); err != nil {
			return false, err
		}
	}

	if valid {
		if valid, err = c.nwItemValidate(ctx, nvItemNetworkKeyDistribute(c.networkKeyDistribute)); err != nil {
			return false, err
		}
	}

	if valid && version.Product != VersionZStack3x0 {
		response, err := c.ZbReadConfiguration(ctx, ZdConfigurationNetworkKey)

		if err != nil {
			return false, err
		}

		valid = bytes.Equal(response.Value, c.networkKey)
	}

	if valid {
		if valid, err = c.nwItemValidate(ctx, nvItemPanID(c.panID)); err != nil {
			return false, err
		}
	}

	if valid {
		if valid, err = c.nwItemValidate(ctx, nvItemExtendedPanID(c.extendedPanID)); err != nil {
			return false, err
		}
	}

	return valid, err
}

func (c *Client) initialization(ctx context.Context) (err error) {
	version, err := c.Version(ctx)
	if err != nil {
		return err
	}

	if err = c.Reset(ctx, ResetTypeSoft); err != nil {
		return err
	}

	// STARTUP_OPTION
	if err = c.nvItemWrite(ctx, nvItemStartupOption()); err != nil {
		return err
	}

	if err = c.Reset(ctx, ResetTypeSoft); err != nil {
		return err
	}

	// logical type as coordinator
	if err = c.nvItemWrite(ctx, nvItemLogicalType(0x00)); err != nil {
		return err
	}

	// network key distribute
	if err = c.nvItemWrite(ctx, nvItemNetworkKeyDistribute(c.networkKeyDistribute)); err != nil {
		return err
	}

	// zdo direct cb
	if err = c.nvItemWrite(ctx, nvItemZdoDirectCb()); err != nil {
		return err
	}

	// channel list
	if err = c.nvItemWrite(ctx, nvItemChannelList(c.channel)); err != nil {
		return err
	}

	// pan id
	if err = c.nvItemWrite(ctx, nvItemPanID(c.panID)); err != nil {
		return err
	}

	// extended pan id
	if err = c.nvItemWrite(ctx, nvItemExtendedPanID(c.extendedPanID)); err != nil {
		return err
	}

	if version.Product == VersionZStack12 {
		// network key
		if err = c.ZbWriteConfiguration(ctx, ZdConfigurationNetworkKey, c.networkKey); err != nil {
			return err
		}

		// TC link key
		if err = c.nvItemWrite(ctx, nvItemTcLinKey()); err != nil {
			return err
		}
	}

	// NV_ITEM_UNINIT
	if err = c.SysOsalNvItemInit(ctx, NvItemIDHasConfiguredZStack1, []byte{0x00}); err != nil {
		return err
	}

	// ZNP has configured
	if err = c.nvItemWrite(ctx, nvItemZnpHasConfigured(version)); err != nil {
		//return err
	}

	c.nwItemValidate(ctx, nvItemZnpHasConfigured(version))

	return nil
}

func (c *Client) defaultWatcher() {
	watcher := c.Watch()
	defer func() {
		c.unregisterWatcher(watcher)
	}()

	for {
		select {
		case frame := <-watcher.NextFrame():
			switch frame.SubSystem() {
			case SubSystemZDOInterface:
				switch frame.CommandID() {
				case CommandTcDeviceInd:
					if msg, err := c.ZDODeviceJoinedMessage(frame); err == nil {
						device := c.Device(msg.NetworkAddress)
						if device != nil {
							device.SetIEEEAddress(msg.ExtendAddress)
						} else {
							d := NewDevice()
							d.SetNetworkAddress(msg.NetworkAddress)
							d.SetIEEEAddress(msg.ExtendAddress)

							c.deviceAdd(d)
						}
					}

				case CommandEndDeviceAnnounceInd:
					if msg, err := c.ZDOEndDeviceAnnounceMessage(frame); err == nil {
						device := c.Device(msg.NetworkAddress)
						if device != nil {
							device.SetCapabilities(msg.Capabilities)
						} else {
							d := NewDevice()
							d.SetNetworkAddress(msg.NetworkAddress)
							d.SetCapabilities(msg.Capabilities)

							c.deviceAdd(d)
						}
					}

				case CommandLeaveInd:
					if msg, err := c.ZDODeviceLeaveMessage(frame); err == nil {
						c.deviceRemove(msg.SourceAddress)
					}

				case CommandPermitJoinInd:
					atomic.StoreUint32(&c.permitJoin, 0)
				}
			}

		case <-watcher.NextError():
		}
	}
}
