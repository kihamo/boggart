package zstack

import (
	"context"
	"errors"
)

/*
AF_REGISTER

This command enables the tester to register an application’s endpoint description.

Usage:

	SREQ:
		          1        |      1      |       1     |     1    |      2    |      2      |     1     |      1     |        1         |       0-32       |         1         |       0-32
		Length = 0x09-0x49 | Cmd0 = 0x24 | Cmd1 = 0x00 | EndPoint | AppProfID | AppDeviceID | AppDevVer | LatencyReq | AppNumInClusters | AppInClusterList | AppNumOutClusters | AppOutClusterList
	Attributes:
		EndPoint          1 byte   Specifies the endpoint of the device
		AppProfID         2 bytes  Specifies the profile Id of the application
		AppDeviceID       2 bytes  Specifies the device description Id for this endpoint
		AddDevVer         1 byte   Specifies the device version number
		LatencyReq        1 byte   Specifies latency.
		                           0x00-No latency
		                           0x01-fast beacons
		                           0x02-slow beacons
		AppNumInClusters  1 byte   The number of Input cluster Id’s following in the AppInClusterList
		AppInClusterList  32 bytes Specifies the list of Input Cluster Id’s
		AppNumOutClusters 1 byte   Specifies the number of Output cluster Id’s following in the AppOutClusterList
		AppOutClusterList 32 bytes Specifies the list of Output Cluster Id’s

	SRSP:
		       1      |       1     |       1     |    1
		Length = 0x01 | Cmd0 = 0x64 | Cmd1 = 0x00 | Status
	Attributes:
		Status 1 byte Status is either Success (0) or Failure (1).

Example from zigbee2mqtt:

	zigbee-herdsman:adapter:zStack:znp:SREQ --> AF - register - {"appdeviceid":5,"appdevver":0,"appnuminclusters":0,"appinclusterlist":[],"appnumoutclusters":0,"appoutclusterlist":[],"latencyreq":0,"endpoint":1,"appprofid":260} +14ms
	zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,9,36,0,1,4,1,5,0,0,0,0,0,44] +14ms
	zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,100,0,0,101] +6ms
	zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,100,0,0,101] +0ms
	zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 3 - 4 - 0 - [0] - 101 +0ms
	zigbee-herdsman:adapter:zStack:znp:SRSP <-- AF - register - {"status":0} +12ms
	zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [] +1ms
*/
func (c *Client) AfRegister(ctx context.Context, endpoint Endpoint) error {
	response, err := c.CallWithResultSREQ(ctx, endpoint.AsBuffer().Frame(0x24, 0x00))
	if err != nil {
		return err
	}

	if response.Command0() != 0x64 {
		return errors.New("bad response")
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() == 0 {
		return errors.New("failure")
	}

	if status := dataOut.ReadCommandStatus(); status != CommandStatusSuccess {
		return status
	}

	return nil
}
