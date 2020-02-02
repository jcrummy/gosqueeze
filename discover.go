package gosqueeze

import (
	"log"
	"net"
)

// Discover returns a list of squeezebox devices found on the network
func Discover(iface *net.Interface) ([]Sb, error) {
	// Put together packet to send
	p := packet{
		DstBroadcast: true,
		DstAddrType:  addrTypeEth,
		DstMac:       macZero,
		SrcBroadcast: false,
		SrcAddrType:  addrTypeUDP,
		SrcIP:        ipZero,
		SrcPort:      0,
		UcpMethod:    ucpMethodAdvDiscover,
	}
	packet := p.assemble()

	var sb []Sb

	err := broadcastReceive(iface, 17784, packet, func(n int, addr *net.UDPAddr, buf []byte) {
		log.Printf("Received %d bytes from %s.\n", n, addr.String())
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodAdvDiscover {
			data, err := parseInboundData(p.Data)
			if err != nil {
				log.Println(err)
				return
			}
			foundSB := Sb{MacAddr: p.SrcMac}
			setDeviceData(&foundSB, data)
			foundSB.GetIP(iface)
			sb = append(sb, foundSB)
		}
	})
	if err != nil {
		return nil, err
	}

	return sb, nil
}
