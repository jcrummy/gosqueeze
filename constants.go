// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package gosqueeze

// Networking constants
var (
	ipZero   = []byte{0x00, 0x00, 0x00, 0x00}
	macZero  = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	portUDAP = []byte{0x45, 0x78}
	portZero = []byte{0x00, 0x00}
)

const (
	udapPort int = 17784
)

// Addressing type of the packet
const (
	addrTypeRaw = iota
	addrTypeEth
	addrTypeUDP
	addrTypeThree
)

// UCP methods
const (
	ucpMethodZero = iota
	ucpMethodDiscover
	ucpMethodGetIP
	ucpMethodSetIP
	ucpMethodReset
	ucpMethodGetData
	ucpMethodSetData
	ucpMethodError
	ucpMethodCredentialsError
	ucpMethodAdvDiscover
	ucpMethodTen
	ucpMethodGetUUID
)

// UCP Codes
const (
	UCPCodeZero = iota
	UCPCodeOne
	UCPCodeDeviceName
	UCPCodeDeviceType
	UCPodeUseDHCP
	UCPCodeIPAddr
	UCPCodeSubnetMask
	UCPCodeGatewayAddr
	UCPCodeEight
	UCPCodeFirmwareRev
	UCPCodeHardwareRev
	UCPCodeDeviceID
	UCPCodeDeviceStatus
	UCPCodeUUID
)

const (
	dataLanIPMode = iota
	dataLanNetworkaddress
	dataLanSubnetMask
	dataLanGateway
	dataLanHostname
	dataLanBridging
	dataLanInterface
	dataLanDNSPrimary
	dataLanDNSSecondary
	dataSqueezeCenterAddress
)

// Misc constants
var (
	uapClassUCP        = []byte{0x00, 0x01, 0x00, 0x01}
	udapTypeUCP        = []byte{0xC0, 0x01}
	defaultCredentials = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // len = 32
)
