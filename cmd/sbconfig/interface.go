package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/c-bata/go-prompt"
)

var ifaces []net.Interface
var i int

func selectedInterface() *net.Interface {
	return &ifaces[i]
}

func selectInterface() {
	var err error
	ifaces, err = net.Interfaces()
	if err != nil {
		panic("Unable to find network interfaces: " + err.Error())
	}
	fmt.Println("The following network interfaces were found:")
	for i, iface := range ifaces {
		fmt.Printf("[%02d] %s - %+v\n", i, iface.Name, iface.HardwareAddr)
	}
	for {
		t := prompt.Input(fmt.Sprintf("Select interface to use [%d]: ", i), func(d prompt.Document) []prompt.Suggest {
			return nil
		})
		i, err = strconv.Atoi(t)
		if err != nil {
			fmt.Println("Invalid selection. Please use a number.")
			continue
		}
		if i >= len(ifaces) {
			fmt.Printf("Invalid selection. Please enter a number from 0 to %d.\n", len(ifaces)-1)
			continue
		}
		break
	}
}
