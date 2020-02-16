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
	Data        deviceData
}

// Tagged as offset,data length (in bytes)
type deviceData struct {
	LanIPMode            bool   `gosqueeze:"4,1"` // Tagged as offset,data length (in bytes)
	LanNetworkAddress    net.IP `gosqueeze:"5,4"`
	LanSubnetMask        net.IP `gosqueeze:"9,4"`
	LanGateway           net.IP `gosqueeze:"13,4"`
	Hostname             string `gosqueeze:"17,33"`
	Bridging             bool   `gosqueeze:"50,1"`
	Interface            uint8  `gosqueeze:"52,1"`
	PrimaryDNS           net.IP `gosqueeze:"59,4"`
	SecondaryDNS         net.IP `gosqueeze:"67,4"`
	ActiveServerAddress  net.IP `gosqueeze:"71,4"`
	SqueezeCenterAddress net.IP `gosqueeze:"79,4"`
	SqueezeCenterName    string `gosqueeze:"83,33"`
	WirelessMode         uint8  `gosqueeze:"173,1"`
	WirelessSSID         string `gosqueeze:"183,33"`
	WirelessChannel      uint8  `gosqueeze:"216,1"`
	WirelessRegion       uint8  `gosqueeze:"218,1"`
	WirelessKeylen       uint8  `gosqueeze:"220,1"`
	WirelessEWPKey0      []byte `gosqueeze:"222,13"`
	WirelessEWPKey1      []byte `gosqueeze:"235,13"`
	WirelessEWPKey2      []byte `gosqueeze:"248,13"`
	WirelessEWPKey3      []byte `gosqueeze:"261,13"`
	WirelessWEPOn        bool   `gosqueeze:"274,1"`
	WirelessWPACipher    uint8  `gosqueeze:"275,1"`
	WirelessWPAMode      uint8  `gosqueeze:"276,1"`
	WirelessWPAOn        bool   `gosqueeze:"277,1"`
	WirelessWPAPSK       string `gosqueeze:"278,64"`
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
