package main

import (
	"fmt"
	"net"

	"github.com/jcrummy/gosqueeze"
)

var sbs map[string]gosqueeze.Sb

func discover(iface *net.Interface) {
	var err error
	sbs, err = gosqueeze.Discover(iface)
	if err != nil {
		fmt.Printf("Error finding devices: %s\n", err.Error())
	}
	fmt.Println("Found the following devices: ")
	for idx, sb := range sbs {
		err = sb.GetIP(iface)
		if err != nil {
			fmt.Printf("Error retrieving IP address: %s\n", err.Error())
		}
		err = sb.GetData(iface)
		if err != nil {
			fmt.Printf("Error retrieving device data: %s\n", err.Error())
		}
		fmt.Printf("  [%02d] %+v at %+v\n", idx, sb.MacAddr, sb.IPAddr)
	}

	return
}
