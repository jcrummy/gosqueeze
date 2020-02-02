package gosqueeze

import (
	"errors"
	"log"
	"net"
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
	Response    []byte
}

// GetIP adds IP address information to the Sb device
func (s *Sb) GetIP(iface *net.Interface) error {
	if s.MacAddr == nil {
		return errors.New("Hardware address required to retrieve IP address")
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

	err := broadcastReceive(iface, 17784, packet, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodGetIP {
			data, err := parseInboundData(p.Data)
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
