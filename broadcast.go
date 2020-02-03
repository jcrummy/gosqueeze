package gosqueeze

import (
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func broadcastReceive(iface *net.Interface, sendPort int, sendMsg []byte, timeout time.Duration,
	handler func(n int, addr *net.UDPAddr, buf []byte)) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(sendPort))
	if err != nil {
		return err
	}
	ifaceAddress, err := getIfaceAddr(iface)
	if err != nil {
		return err
	}
	sendAddr, err := net.ResolveUDPAddr("udp", ifaceAddress+":0")
	if err != nil {
		return err
	}

	// send packet
	sconn, err := net.ListenUDP("udp", sendAddr)
	if err != nil {
		return err
	}
	// We listen on the same port chosen to send on, so we need to figure it out now
	listenPort := strings.Split(sconn.LocalAddr().String(), ":")[1]
	listenAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+listenPort)
	if err != nil {
		return err
	}
	_, err = sconn.WriteToUDP(sendMsg, broadcastAddr)
	if err != nil {
		return err
	}
	sconn.Close()

	// Setup listener on all interfaces because we can't listen to broadcast messages on
	// a specific interface
	rconn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	defer rconn.Close()

	buf := make([]byte, 1024)
	rconn.SetDeadline(time.Now().Add(timeout))
	for {
		n, addr, err := rconn.ReadFromUDP(buf)
		if err != nil {
			if err.(net.Error).Timeout() {
				break
			} else {
				log.Println(err)
			}
		}
		handler(n, addr, buf)
	}

	return nil
}

func broadcastSingle(iface *net.Interface, sendPort int, sendMsg []byte, timeout time.Duration,
	handler func(n int, addr *net.UDPAddr, buf []byte)) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+strconv.Itoa(sendPort))
	if err != nil {
		return err
	}
	ifaceAddress, err := getIfaceAddr(iface)
	if err != nil {
		return err
	}
	sendAddr, err := net.ResolveUDPAddr("udp", ifaceAddress+":0")
	if err != nil {
		return err
	}

	// send packet
	sconn, err := net.ListenUDP("udp", sendAddr)
	if err != nil {
		return err
	}
	// We listen on the same port chosen to send on, so we need to figure it out now
	listenPort := strings.Split(sconn.LocalAddr().String(), ":")[1]
	listenAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+listenPort)
	if err != nil {
		return err
	}
	_, err = sconn.WriteToUDP(sendMsg, broadcastAddr)
	if err != nil {
		return err
	}
	sconn.Close()

	// Setup listener on all interfaces because we can't listen to broadcast messages on
	// a specific interface
	rconn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	defer rconn.Close()

	buf := make([]byte, 1024)
	rconn.SetDeadline(time.Now().Add(timeout))
	n, addr, err := rconn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	handler(n, addr, buf)

	return nil
}

func getIfaceAddr(iface *net.Interface) (string, error) {
	laddrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	var ifaceAddress string
	for _, addr := range laddrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			a := ipnet.IP.To4()
			if a == nil {
				continue
			}
			ifaceAddress = a.String()
		}
	}
	if ifaceAddress == "" {
		return "", errors.New("No addresses associated with interface")
	}
	return ifaceAddress, nil
}
