package gosqueeze

import (
	"encoding/binary"
	"errors"
	"fmt"
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
	Data        deviceData
}

// Tagged as offset,data length (in bytes)
type deviceData struct {
	LanIPMode            bool   `gosqueeze:"4,1"`    // false = static IP, true = DHCP
	LanNetworkAddress    net.IP `gosqueeze:"5,4"`    // static network address
	LanSubnetMask        net.IP `gosqueeze:"9,4"`    // static subnet mask
	LanGateway           net.IP `gosqueeze:"13,4"`   // static gateway address
	Hostname             string `gosqueeze:"17,33"`  // device hostname
	Bridging             bool   `gosqueeze:"50,1"`   // true = use device as wireless bridge
	Interface            uint8  `gosqueeze:"52,1"`   // 0 = use Wireless, 1 = use Wired
	PrimaryDNS           net.IP `gosqueeze:"59,4"`   // static primary DNS address
	SecondaryDNS         net.IP `gosqueeze:"67,4"`   // static secondary DNS address
	ActiveServerAddress  net.IP `gosqueeze:"71,4"`   // IP address of currently active server (read only)
	SqueezeCenterAddress net.IP `gosqueeze:"79,4"`   // IP address of local Squeezecenter server
	SqueezeCenterName    string `gosqueeze:"83,33"`  // Name of local Squeezecenter server (read only)
	WirelessMode         uint8  `gosqueeze:"173,1"`  // 0 = infrastructure, 1 = Ad Hoc
	WirelessSSID         string `gosqueeze:"183,33"` // SSID of WiFi access point to connect to
	WirelessChannel      uint8  `gosqueeze:"216,1"`  // WiFi Channel, can normally leave at 0
	WirelessRegion       uint8  `gosqueeze:"218,1"`  // 4 = US, 6 = CA, 7 = AU, 13 = FR, 14 = EU, 16 = JP, 21 = TW, 23 = CH
	WirelessKeylen       uint8  `gosqueeze:"220,1"`  // Length of wireless key (0 = 64-bit, 1 = 128-bit)
	WirelessEWPKey0      []byte `gosqueeze:"222,13"` // WEP key 0 - in Hex
	WirelessEWPKey1      []byte `gosqueeze:"235,13"` // WEP key 1 - in Hex
	WirelessEWPKey2      []byte `gosqueeze:"248,13"` // WEP key 2 - in Hex
	WirelessEWPKey3      []byte `gosqueeze:"261,13"` // WEP key 3 - in Hex
	WirelessWEPOn        bool   `gosqueeze:"274,1"`  // 0 = Wep Off, 1 = Wep On
	WirelessWPACipher    uint8  `gosqueeze:"275,1"`  // 1 = TKIIP, 2 = AES, 3 = TKIP & AES
	WirelessWPAMode      uint8  `gosqueeze:"276,1"`  // 1 = WPA, 2 = WPA2
	WirelessWPAOn        bool   `gosqueeze:"277,1"`  // 0 = WPA Off, 1 = WPA On
	WirelessWPAPSK       string `gosqueeze:"278,64"` // WPA Public Shared Key
}

// GetIP retrieves IP address information from the SqueezeBox device
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
			s.populateFields(data)
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

// GetData retrieves all data points from the SqueezeBox device
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
	p.setDataRetrieve(s.Data)
	packet := p.assemble()

	err := broadcastSingle(iface, udapPort, packet, 500*time.Millisecond, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodGetData {
			err = p.parseData(&s.Data)
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

// SaveData saves all current values to the SqueezeBox device
func (s *Sb) SaveData(iface *net.Interface) error {
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
	}
	numDataFields := p.setDataForSave(s.Data)
	packet := p.assemble()

	err := broadcastSingle(iface, udapPort, packet, 500*time.Millisecond, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodSetData {
			numberChanged := int(binary.BigEndian.Uint16(p.Data))
			if numberChanged == numDataFields {
				fmt.Println("Successfully set data.")
			} else {
				fmt.Println("Error setting data. Not all fields were saved")
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// populateFields sets the Sb root field values based on the
// provided map.
func (s *Sb) populateFields(f fields) {
	for i, v := range f {
		switch i {
		//case UCPCodeZero:
		//case UCPCodeOne:
		case UCPCodeDeviceName:
			s.Name = string(v)
		case UCPCodeDeviceType:
			s.Type = string(v)
		// case UCPodeUseDHCP    :
		case UCPCodeIPAddr:
			s.IPAddr = v
		case UCPCodeSubnetMask:
			s.SubnetMask = v
		case UCPCodeGatewayAddr:
			s.GatewayAddr = v
		// case UCPCodeEight       :
		case UCPCodeFirmwareRev:
			s.FirmwareRev = uint(binary.BigEndian.Uint16(v))
		case UCPCodeHardwareRev:
			s.HardwareRev = uint(binary.BigEndian.Uint32(v))
		case UCPCodeDeviceID:
			s.ID = uint(binary.BigEndian.Uint16(v))
		case UCPCodeDeviceStatus:
			s.Status = string(v)
			// case UCPCodeUUID :
		}
	}
}
