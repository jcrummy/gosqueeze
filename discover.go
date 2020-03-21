// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package gosqueeze

import (
	"log"
	"net"
	"time"

	"github.com/jcrummy/gosqueeze/internal/broadcast"
	"github.com/jcrummy/gosqueeze/internal/constants"
	"github.com/jcrummy/gosqueeze/internal/packet"
)

// Discover returns a list of squeezebox devices found on the network
func Discover(iface *net.Interface) ([]Sb, error) {
	// Put together packet to send
	p := packet.Packet{
		DstBroadcast: true,
		DstAddrType:  constants.AddrTypeEth,
		DstMac:       constants.MacZero,
		SrcBroadcast: false,
		SrcAddrType:  constants.AddrTypeUDP,
		SrcIP:        constants.IPZero,
		SrcPort:      0,
		UcpMethod:    constants.UCPMethodAdvDiscover,
	}
	packetBytes := p.Assemble()

	var sb []Sb

	err := broadcast.BroadcastReceive(iface, 17784, packetBytes, 3*time.Second, func(n int, addr *net.UDPAddr, buf []byte) {
		p, err := packet.Parse(buf[:n])
		if err != nil {
			return
		}
		if p.UcpMethod == constants.UCPMethodAdvDiscover {
			data, err := p.ParseFields()
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
