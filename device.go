package gosqueeze

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

// Sb represents a squeezebox receiver device
type Sb struct {
	MacAddr     net.HardwareAddr
	IPAddr      net.IP
	GatewayAddr net.IP
	SubnetMask  net.IPMask
	ID          uint
	Type        string
	Name        string
	Status      string
	HardwareRev uint
	FirmwareRev uint
	Data        map[int][]byte
}

type deviceData struct {
	LanIPMode         bool   `gosqueeze:"4,1"`
	LanNetworkAddress net.IP `gosqueeze:"5,4"`
	LanSubnetMask     net.IP `gosqueeze:"9,4"`
	// 13:  dataItem{"lanGateway", "IP address of default network gateway", 4},
	// 17:  dataItem{"hostname", "Device hostname (is this set automatically?)", 33},
	// 50:  dataItem{"bridging", "Use device as a wireless bridge (not sure about this)", 1},
	// 52:  dataItem{"interface", "0 - wireless, 1 - wired (is set to 128 after factory reset)", 1},
	// 59:  dataItem{"primaryDNS", "IP address of primary DNS server", 4},
	// 67:  dataItem{"secondaryDNS", "IP address of secondary DNS server", 4},
	// 71:  dataItem{"activeServerAddress", "IP address of currently active server (either Squeezenetwork or local server)", 4},
	// 79:  dataItem{"squeezeCenterAddress", "IP address of local SqueezeCenter server", 4},
	// 83:  dataItem{"squeezeCenterName", "Name of local SqueezeCenter server (???)", 33},
	// 173: dataItem{"wirelessMode", "0 - Infrastructure, 1 - Ad Hoc", 1},
	// 183: dataItem{"wirelessSSID", "Wireless network name", 33},
	// 216: dataItem{"wirelessChannel", "Wireless channel (used by AdHoc mode???)", 1},
	// 218: dataItem{"wirelessRegion", "4 - US, 6 - CA, 7 - AU, 13 - FR, 14 - EU, 16 - JP, 21 - TW, 23 - CH", 1},
	// 220: dataItem{"wirelessKeylen", "Length of wireless key, (0 - 64-bit, 1 - 128-bit)", 1},
	// 222: dataItem{"wirelessWEPKey0", "WEP key 0 - enter in hex", 13},
	// 235: dataItem{"wirelessWEPKey1", "WEP key 0 - enter in hex", 13},
	// 248: dataItem{"wirelessWEPKey2", "WEP key 0 - enter in hex", 13},
	// 261: dataItem{"wirelessWEPKey3", "WEP key 0 - enter in hex", 13},
	// 274: dataItem{"wirelessWEPOn", "0 - WEP Off, 1 - WEP On", 1},
	// 275: dataItem{"wirelessWPACipher", "1 - TKIP, 2 - AES, 3 - TKIP & AES", 1},
	// 276: dataItem{"wirelessWPAMode", "1 - WPA, 2 - WPA2", 1},
	// 277: dataItem{"wirelessWPAOn", "0 - WPA Off, 1 - WPA On", 1},
	// 278: dataItem{"wirelessSPAPSK", "WPA Public Shared Key", 64},
}

// GetIP adds IP address information to the Sb device
func (s *Sb) GetIP(iface *net.Interface) error {
	if s.MacAddr == nil {
		return errors.New("Hardware address required")
	}

	p := packet{
		DstBroadcast: false,
		DstAddrType:  addrTypeEth,
		DstMac:       s.MacAddr,
		SrcBroadcast: false,
		SrcAddrType:  addrTypeUDP,
		SrcIP:        ipZero,
		SrcPort:      0,
		UcpMethod:    ucpMethodGetIP,
	}
	packet := p.assemble()

	err := broadcastSingle(iface, 17784, packet, 500*time.Millisecond, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodGetIP {
			data, err := p.parseFields()
			if err != nil {
				log.Println(err)
				return
			}
			setDeviceData(s, data)
		}
	})
	if err != nil {
		return err
	}
	if s.IPAddr == nil {
		return errors.New("Error retrieving IP address")
	}
	return nil
}

// GetData adds all information to the Sb device
func (s *Sb) GetData(iface *net.Interface) error {
	if s.MacAddr == nil {
		return errors.New("Hardware address required")
	}

	p := packet{
		DstBroadcast: false,
		DstAddrType:  addrTypeEth,
		DstMac:       s.MacAddr,
		SrcBroadcast: false,
		SrcAddrType:  addrTypeUDP,
		SrcIP:        ipZero,
		SrcPort:      0,
		UcpMethod:    ucpMethodGetData,
	}
	packet := p.assemble()

	err := broadcastSingle(iface, udapPort, packet, 500*time.Millisecond, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodGetData {
			s.Data, err = p.parseData()
			if err != nil {
				log.Println("Error getting data from device: " + err.Error())
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// SetData sets to provided data on the device
func (s *Sb) SetData(iface *net.Interface) error {
	if s.MacAddr == nil {
		return errors.New("Hardware address required")
	}

	p := packet{
		DstBroadcast: false,
		DstAddrType:  addrTypeEth,
		DstMac:       s.MacAddr,
		SrcBroadcast: false,
		SrcAddrType:  addrTypeUDP,
		SrcIP:        ipZero,
		SrcPort:      0,
		UcpMethod:    ucpMethodSetData,
		Data:         []byte{0x00, 0x01, 0x00, 0x05, 0x00, 0x04, 192, 168, 18, 25},
	}
	packet := p.assemble()

	err := broadcastSingle(iface, udapPort, packet, 500*time.Millisecond, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodSetData {
			numberChanged := binary.BigEndian.Uint16(p.Data)
			if numberChanged == 1 {
				log.Println("Successfully set data.")
			} else {
				log.Println("Error setting data.")
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}
