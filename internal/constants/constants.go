// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package constants

// Networking constants
var (
	IPZero   = []byte{0x00, 0x00, 0x00, 0x00}
	MacZero  = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	PortUDAP = []byte{0x45, 0x78}
	PortZero = []byte{0x00, 0x00}
)

// UdapPort is the standard communication port
const (
	UdapPort int = 17784
)

// Addressing type of the packet
const (
	AddrTypeRaw = iota
	AddrTypeEth
	AddrTypeUDP
	AddrTypeThree
)

// UCP methods
const (
	UCPMethodZero = iota
	UCPMethodDiscover
	UCPMethodGetIP
	UCPMethodSetIP
	UCPMethodReset
	UCPMethodGetData
	UCPMethodSetData
	UCPMethodError
	UCPMethodCredentialsError
	UCPMethodAdvDiscover
	UCPMethodTen
	UCPMethodGetUUID
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
	DataLanIPMode = iota
	DataLanNetworkaddress
	DataLanSubnetMask
	DataLanGateway
	DataLanHostname
	DataLanBridging
	DataLanInterface
	DataLanDNSPrimary
	DataLanDNSSecondary
	DataSqueezeCenterAddress
)

// Misc constants
var (
	UapClassUCP        = []byte{0x00, 0x01, 0x00, 0x01}
	UdapTypeUCP        = []byte{0xC0, 0x01}
	DefaultCredentials = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} // len = 32
)
