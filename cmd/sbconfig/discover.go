package main

import (
	"fmt"
	"net"

	"github.com/jcrummy/gosqueeze"
)

var sbs []gosqueeze.Sb

func discover(iface *net.Interface) {
	var err error
	sbs, err = gosqueeze.Discover(iface)
	if err != nil {
		fmt.Printf("Error finding devices: %s\n", err.Error())
	}
	fmt.Println("Found the following devices: ")
	for i := 0; i < len(sbs); i++ {
		err = sbs[i].GetIP(iface)
		if err != nil {
			fmt.Printf("Error retrieving IP address: %s\n", err.Error())
		}
		err = sbs[i].GetData(iface)
		if err != nil {
			fmt.Printf("Error retrieving device data: %s\n", err.Error())
		}
		fmt.Printf("  [%02d] %+v at %+v\n", i, sbs[i].MacAddr, sbs[i].IPAddr)
	}

	return
}
