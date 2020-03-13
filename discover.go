// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package gosqueeze

import (
	"log"
	"net"
	"time"
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

	err := broadcastReceive(iface, 17784, packet, 3*time.Second, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := parsePacket(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == ucpMethodAdvDiscover {
			data, err := p.parseFields()
			if err != nil {
				log.Println(err)
				return
			}
			foundSB := Sb{MacAddr: p.SrcMac}
			foundSB.populateFields(data)
			sb = append(sb, foundSB)
		}
	})
	if err != nil {
		return nil, err
	}

	return sb, nil
}
