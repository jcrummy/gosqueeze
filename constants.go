package gosqueeze

const udpMaxMsgLen int = 1500

// Networking constants
var (
	dstTypeEth = []byte{0x01}
	ipZero     = []byte{0x00, 0x00, 0x00, 0x00}
	macZero    = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	portUDAP   = []byte{0x45, 0x78}
	portZero   = []byte{0x00, 0x00}
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

// Misc constants
var (
	uapClassUCP = []byte{0x00, 0x01, 0x00, 0x01}
	udapTypeUCP = []byte{0xC0, 0x01}
)
