// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package gosqueeze provides an interface to configure Logitech
SqueezeBox devices over a network.

This package is inspired by the work of Robin Bowes and his Net-UDAP
Perl module. In particular the packet format was deciphered based on his
work.

This module was created specifically for use in the github.com/jcrummy/sbconfig
program, however it is available for use in other contexts as well.

Basics

You must specificy a network interface to use for sending broadcast messages. Note
replies are listened for on all interfaces due to limitations in broadcast handling.

	// Discover what devices are available on the network
	sbs, _ = gosqueeze.Discover(iface)

	// Get information about a specific device
	sbs[0].GetIP(iface)
	sbs[0].GetData(iface)

	// Modify a configuration point on the device
	sbs[0].Data.WirelessSSID = "NewSSID"

	// Save configuration changes to the device
	sbs[0].SaveData(iface)

*/
package gosqueeze
