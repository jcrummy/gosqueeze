package sb

import (
	"fmt"
	"net"
	"strings"

	"github.com/jcrummy/gosqueeze"

	"github.com/c-bata/go-prompt"
)

// Configure opens prompt to configure a specific device
func Configure(device *gosqueeze.Sb, iface *net.Interface) {
	c := configurator{
		device: device,
		iface:  iface,
	}
	fmt.Println("Use Ctrl-D to exit configuration mode.")
	t := prompt.New(c.executor, configureCompleter,
		prompt.OptionTitle(fmt.Sprintf("sbconfig: Configuring %+v", c.device.MacAddr)),
		prompt.OptionPrefix("C>> "),
	)
	t.Run()
	return
}

type configurator struct {
	device *gosqueeze.Sb
	iface  *net.Interface
}

func configureCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "show", Description: "Show current settings"},
		//{Text: "exit", Description: "Exit program"},
		{Text: "set", Description: "Set a particular value"},
		{Text: "save", Description: "Save current values to device"},
	}
	setpoints := []prompt.Suggest{
		{Text: "LanIPMode", Description: "False = Static IP, True = DHCP"},
		{Text: "LanNetworkAddress", Description: "Static IP address"},
		{Text: "LanSubnetMask", Description: "Static subnet mask"},
		{Text: "LanGateway", Description: "Static gateway address"},
		{Text: "Hostname", Description: "Device hostname"},
		{Text: "Bridging", Description: "True = use device as wireless bridge"},
		{Text: "Interface", Description: "0 = Use wireless link, 1 = Use wired link"},
		{Text: "PrimaryDNS", Description: "Static primary DNS address"},
		{Text: "SecondaryDNS", Description: "Static secondary DNS address"},
		{Text: "SqueezeCenterAddress", Description: "IP address of local Squeezecenter server"},
		{Text: "Wireless Mode", Description: "0 = Infrastructure, 1 = Ad Hoc"},
		{Text: "WirelessSSID", Description: "SSID of WiFi access point to connect to"},
		{Text: "WirelessChannel", Description: "WiFi channel, 0 for automatic"},
		{Text: "WirelessRegion", Description: "4 US, 6 CA, 7 AU, 13 FR, 14 EU, 16 JP, 21 TW, 23 CH"},
		{Text: "WirelessKeylen", Description: "Length of wireless key (0 = 64-bit, 1 = 128-bit)"},
		{Text: "WirelessWEPKey0", Description: "WEP key 0 - in hex"},
		{Text: "WirelessWEPKey1", Description: "WEP key 1 - in hex"},
		{Text: "WirelessWEPKey2", Description: "WEP key 2 - in hex"},
		{Text: "WirelessWEPKey3", Description: "WEP key 3 - in hex"},
		{Text: "WirelessWEPOn", Description: "0 = Wep Off, 1 = Wep On"},
		{Text: "WirelessWPACipher", Description: "1 = TKIP, 2 = AES, 3 = TKIP & AES"},
		{Text: "WirelessWPAMode", Description: "1 = WPA, 2 = WPA2"},
		{Text: "WirelessWPAOn", Description: "0 = WPA Off, 1 = WPA On"},
		{Text: "WirelessWPAPSK", Description: "WPA Public Shared Key"},
	}
	cmds := strings.Split(d.Text, " ")
	if cmds[0] == "set" {
		if len(cmds) > 1 {
			return prompt.FilterHasPrefix(setpoints, cmds[1], true)
		}
		return prompt.FilterHasPrefix(setpoints, d.GetWordBeforeCursorUntilSeparator(" "), true)
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (c *configurator) showValues() {
	fmt.Printf("LanIPMode: %+v\n", c.device.Data.LanIPMode)
	fmt.Printf("LanNetworkAddress: %+v\n", c.device.Data.LanNetworkAddress)
	fmt.Printf("LanSubnetMask: %+v\n", c.device.Data.LanSubnetMask)
	fmt.Printf("LanGateway: %+v\n", c.device.Data.LanGateway)
	fmt.Printf("Hostname: %+v\n", c.device.Data.Hostname)
	fmt.Printf("Bridging: %+v\n", c.device.Data.Bridging)
	fmt.Printf("Interface: %+v\n", c.device.Data.Interface)
	fmt.Printf("PrimaryDNS: %+v\n", c.device.Data.PrimaryDNS)
	fmt.Printf("SecondaryDNS: %+v\n", c.device.Data.SecondaryDNS)
	fmt.Printf("ActiveServerAddress: %+v\n", c.device.Data.ActiveServerAddress)
	fmt.Printf("SqueezeCenterAddress: %+v\n", c.device.Data.SqueezeCenterAddress)
	fmt.Printf("SqueezeCenterName: %+v\n", c.device.Data.SqueezeCenterName)
	fmt.Printf("WirelessMode: %+v\n", c.device.Data.WirelessMode)
	fmt.Printf("WirelessSSID: %+v\n", c.device.Data.WirelessSSID)
	fmt.Printf("WirelessChannel: %+v\n", c.device.Data.WirelessChannel)
	fmt.Printf("WirelessRegion: %+v\n", c.device.Data.WirelessRegion)
	fmt.Printf("WirelessKeylen: %+v\n", c.device.Data.WirelessKeylen)
	fmt.Printf("WirelessWEPKey0: %+v\n", c.device.Data.WirelessWEPKey0)
	fmt.Printf("WirelessWEPKey1: %+v\n", c.device.Data.WirelessWEPKey1)
	fmt.Printf("WirelessWEPKey2: %+v\n", c.device.Data.WirelessWEPKey2)
	fmt.Printf("WirelessWEPKey3: %+v\n", c.device.Data.WirelessWEPKey3)
	fmt.Printf("WirelessWEPOn: %+v\n", c.device.Data.WirelessWEPOn)
	fmt.Printf("WirelessWPACipher: %+v\n", c.device.Data.WirelessWPACipher)
	fmt.Printf("WirelessWPAMode: %+v\n", c.device.Data.WirelessWPAMode)
	fmt.Printf("WirelessWPAOn: %+v\n", c.device.Data.WirelessWPAOn)
	fmt.Printf("WirelessWPAPSK: %+v\n", c.device.Data.WirelessWPAPSK)
}

func (c *configurator) setValue(s string) {
	vals := strings.Split(s, " ")
	if len(vals) < 3 {
		fmt.Println("set requires a field name and a value.")
		return
	}
	switch strings.ToLower(vals[1]) {
	case "lanipmode":
		setBool(vals[2], &c.device.Data.LanIPMode)

	case "lannetworkaddress":
		setIPAddress(vals[2], &c.device.Data.LanNetworkAddress)

	case "lansubnetmask":
		setIPAddress(vals[2], &c.device.Data.LanSubnetMask)

	case "langateway":
		setIPAddress(vals[2], &c.device.Data.LanGateway)

	case "hostname":
		c.device.Data.Hostname = vals[2]

	case "bridging":
		setBool(vals[2], &c.device.Data.Bridging)

	case "interface":
		setUint8(vals[2], &c.device.Data.Interface, maxUint8(1))

	case "primarydns":
		setIPAddress(vals[2], &c.device.Data.PrimaryDNS)

	case "secondarydns":
		setIPAddress(vals[2], &c.device.Data.SecondaryDNS)

	case "squeezecenteraddress":
		setIPAddress(vals[2], &c.device.Data.SqueezeCenterAddress)

	case "wirelessmode":
		setUint8(vals[2], &c.device.Data.WirelessMode, maxUint8(1))

	case "wirelessssid":
		c.device.Data.WirelessSSID = strings.TrimPrefix(s, vals[0]+" "+vals[1]+" ")

	case "wirelesschannel":
		setUint8(vals[2], &c.device.Data.WirelessChannel, anyUint8())

	case "wirelessregion":
		setUint8(vals[2], &c.device.Data.WirelessRegion, inSetUint8([]uint8{4, 6, 7, 13, 14, 16, 21, 23}))

	case "wirelesskeylen":
		setUint8(vals[2], &c.device.Data.WirelessKeylen, maxUint8(1))

	case "wirelesswepkey0":
		c.device.Data.WirelessWEPKey0 = []byte(vals[2])

	case "wirelesswepkey1":
		c.device.Data.WirelessWEPKey1 = []byte(vals[2])

	case "wirelesswepkey2":
		c.device.Data.WirelessWEPKey2 = []byte(vals[2])

	case "wirelesswepkey3":
		c.device.Data.WirelessWEPKey3 = []byte(vals[2])

	case "wirelesswepon":
		setBool(vals[2], &c.device.Data.WirelessWEPOn)

	case "wirelesswpacipher":
		setUint8(vals[2], &c.device.Data.WirelessWPACipher, maxUint8(3))

	case "wirelesswpamode":
		setUint8(vals[2], &c.device.Data.WirelessWPAMode, maxUint8(2))

	case "wirelesswpaon":
		setBool(vals[2], &c.device.Data.WirelessWPAOn)

	case "wirelesswpapsk":
		c.device.Data.WirelessWPAPSK = strings.TrimPrefix(s, vals[0]+" "+vals[1]+" ")

	default:
		fmt.Println("Invalid data point specified")
	}
}

func (c *configurator) saveValues() {
	c.device.SaveData(c.iface)
}
