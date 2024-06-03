// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package gosqueeze

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jcrummy/gosqueeze/internal/broadcast"
	"github.com/jcrummy/gosqueeze/internal/constants"
	"github.com/jcrummy/gosqueeze/internal/packet"
)

// Discover returns a list of squeezebox devices found on the network
func Discover(iface *net.Interface) (map[string]Sb, error) {
	// Put together packet to send
	p_send := packet.Packet{
		DstBroadcast: true,
		DstAddrType:  constants.AddrTypeEth,
		DstMac:       constants.MacZero,
		SrcBroadcast: false,
		SrcAddrType:  constants.AddrTypeUDP,
		SrcIP:        constants.IPZero,
		SrcPort:      0,
		UcpMethod:    constants.UCPMethodAdvDiscover,
	}
	packetBytes := p_send.Assemble()

	sb := make(map[string]Sb)

	i := 0

	err := broadcast.BroadcastReceive(iface, 17784, packetBytes, 3*time.Second, func(n int, addr *net.UDPAddr, buf []byte) {
		i++
		fmt.Printf("Loop count is %d\n", i)
		p, err := packet.Parse(buf[:n])
		if err != nil {
			return
		}
		//fmt.Println("----------------------------------------------------------------")
		//fmt.Printf("%d\n%+v\n\n", n, p)
		if p.UcpMethod == constants.UCPMethodAdvDiscover {
			data, err := p.ParseFields()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println("----------------------------------------------------------------")
			fmt.Printf("%+v\n\n", data)
			foundSB := Sb{MacAddr: p.SrcMac}
			// fmt.Println("----------------------------------------------------------------")
			// fmt.Printf("%+v\n\n", foundSB)
			foundSB.populateFields(data)
			// fmt.Println("----------------------------------------------------------------")
			fmt.Printf("Memory address of foundSB: %p\n%+v\n\n", &foundSB, foundSB)
			sb[foundSB.MacAddr.String()] = foundSB
			fmt.Printf("%+v\n", sb)
			fmt.Println("----------------------------------------------------------------")
		}
	})
	if err != nil {
		return nil, err
	}

	return sb, nil
}
