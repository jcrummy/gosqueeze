package sb

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func setIPAddress(val string, target *net.IP) {
	ip := net.ParseIP(val)
	if ip == nil {
		fmt.Println("Invalid IP address - write in form of x.x.x.x")
		return
	}
	*target = ip.To4()
}

func setBool(val string, target *bool) {
	switch strings.ToLower(val) {
	case "true":
		*target = true
	case "yes":
		*target = true
	case "1":
		*target = true
	case "false":
		*target = false
	case "no":
		*target = false
	case "0":
		*target = false
	default:
		fmt.Println("Invalid true/false value - value not changed")
	}
}

func setUint8(val string, target *uint8, validNum func(uint8) bool) {
	u, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Not a number - value not changed")
		return
	}
	if !validNum(uint8(u)) {
		fmt.Println("Invalid number - value not changed")
		return
	}
	*target = uint8(u)
}

func maxUint8(max uint8) func(uint8) bool {
	return func(u uint8) bool {
		if u > max {
			return false
		}
		return true
	}
}

func anyUint8() func(uint8) bool {
	return func(u uint8) bool {
		return true
	}
}

func inSetUint8(set []uint8) func(uint8) bool {
	return func(u uint8) bool {
		for _, v := range set {
			if v == u {
				return true
			}
		}
		return false
	}
}
