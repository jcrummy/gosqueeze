---
# Feel free to add content and custom Front Matter to this file.
# To modify the layout, see https://jekyllrb.com/docs/themes/#overriding-theme-defaults

layout: default
---

Usage
=====

1. Download the aproporiate executable using the links above.

2. Set the Squeezebox into configuration mode.

   For Squeezebox Receivee see [SBR front button and LED](http://wiki.slimdevices.com/index.php/SBR_front_button_and_LED) for instructions.

3. Run sbconfig

```
$ ./sbconfig-linux-amd64
The following network interfaces were found:
[00] lo - 
[01] enp5s0f1 - xx:xx:xx:xx:xx:xx
[02] wlp4s0 - xx:xx:xx:xx:xx:xx
Select interface to use [0]: 1
Found the following devices: 
  [00] 00:04:20:xx:xx:xx at 0.0.0.0
>>> configure 0
Configuring device #0 (00:04:20:xx:xx:xx).
Use Ctrl-D to exit configuration mode.
C>> show
LanIPMode: true
LanNetworkAddress: 0.0.0.0
LanSubnetMask: 255.255.255.0
LanGateway: 192.168.18.1
Hostname: TestSqueezer
Bridging: false
Interface: 0
PrimaryDNS: 0.0.0.0
SecondaryDNS: 0.0.0.0
ActiveServerAddress: 0.0.0.0
SqueezeCenterAddress: 192.168.12.3
SqueezeCenterName: 
WirelessMode: 0
WirelessSSID: apSSID
WirelessChannel: 0
WirelessRegion: 6
WirelessKeylen: 0
WirelessEWPKey0: []
WirelessEWPKey1: []
WirelessEWPKey2: []
WirelessEWPKey3: []
WirelessWEPOn: false
WirelessWPACipher: 3
WirelessWPAMode: 2
WirelessWPAOn: true
WirelessWPAPSK: securePassword
C>> set WirelessSSID newSSID
C>> save
Successfully set data.
C>> 
>>> exit
Bye!
$ 
```