package gosqueeze

const udpMaxMsgLen int = 1500

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

// Data Codes
type dataItem struct {
	name   string
	help   string
	length int
}

var dataItems = map[int]dataItem{
	4:   dataItem{"lanIPMode", "0 - Use static IP details, 1 - use DHCP", 1},
	5:   dataItem{"lanNetworkAddress", "IP address of device", 4},
	9:   dataItem{"lanSubnetMask", "Subnet mask of local network", 4},
	13:  dataItem{"lanGateway", "IP address of default network gateway", 4},
	17:  dataItem{"hostname", "Device hostname (is this set automatically?)", 33},
	50:  dataItem{"bridging", "Use device as a wireless bridge (not sure about this)", 1},
	52:  dataItem{"interface", "0 - wireless, 1 - wired (is set to 128 after factory reset)", 1},
	59:  dataItem{"primaryDNS", "IP address of primary DNS server", 4},
	67:  dataItem{"secondaryDNS", "IP address of secondary DNS server", 4},
	71:  dataItem{"activeServerAddress", "IP address of currently active server (either Squeezenetwork or local server)", 4},
	79:  dataItem{"squeezeCenterAddress", "IP address of local SqueezeCenter server", 4},
	83:  dataItem{"squeezeCenterName", "Name of local SqueezeCenter server (???)", 33},
	173: dataItem{"wirelessMode", "0 - Infrastructure, 1 - Ad Hoc", 1},
	183: dataItem{"wirelessSSID", "Wireless network name", 33},
	216: dataItem{"wirelessChannel", "Wireless channel (used by AdHoc mode???)", 1},
	218: dataItem{"wirelessRegion", "4 - US, 6 - CA, 7 - AU, 13 - FR, 14 - EU, 16 - JP, 21 - TW, 23 - CH", 1},
	220: dataItem{"wirelessKeylen", "Length of wireless key, (0 - 64-bit, 1 - 128-bit)", 1},
	222: dataItem{"wirelessWEPKey0", "WEP key 0 - enter in hex", 13},
	235: dataItem{"wirelessWEPKey1", "WEP key 0 - enter in hex", 13},
	248: dataItem{"wirelessWEPKey2", "WEP key 0 - enter in hex", 13},
	261: dataItem{"wirelessWEPKey3", "WEP key 0 - enter in hex", 13},
	274: dataItem{"wirelessWEPOn", "0 - WEP Off, 1 - WEP On", 1},
	275: dataItem{"wirelessWPACipher", "1 - TKIP, 2 - AES, 3 - TKIP & AES", 1},
	276: dataItem{"wirelessWPAMode", "1 - WPA, 2 - WPA2", 1},
	277: dataItem{"wirelessWPAOn", "0 - WPA Off, 1 - WPA On", 1},
	278: dataItem{"wirelessSPAPSK", "WPA Public Shared Key", 64},
}

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
