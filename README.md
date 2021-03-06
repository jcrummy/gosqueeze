Gosqueeze - Squeeze Box Controller
==================================

This package is inspired by the work of Robin Bowes and his Net-UDAP
Perl module. In particular the packet format was deciphered based on his
work.

This module was created specifically for use in the github.com/jcrummy/sbconfig
program, however it is available for use in other contexts as well.

Basics
------

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


Getting the squeeze box setup
-----------------------------
From wiki.slimdevices.com/index.php/SBRFrontButtonAndLED.

To go to setup mode:
1. Press and hold the button for about 3 seconds or until it blinks slow *red* then release it.
2. The LED will go red solid, which means it is booting up. This will take a couple of seconds.
3. The LED will start slowly blinking red, which means it is in setup mode.