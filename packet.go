package gosqueeze

import (
	"encoding/binary"
	"errors"
	"net"
)

type packet struct {
	DstBroadcast bool
	DstAddrType  int
	DstMac       net.HardwareAddr
	DstIP        net.IP
	DstPort      uint
	SrcBroadcast bool
	SrcAddrType  int
	SrcMac       net.HardwareAddr
	SrcIP        net.IP
	SrcPort      uint
	Seq          int
	UdapType     int
	UcpFlags     byte
	UapClass     []byte
	UcpMethod    int
	Data         []byte
}

func parsePacket(buf []byte) (*packet, error) {
	if len(buf) < 27 {
		return nil, errors.New("Packet length too short")
	}

	var p packet
	i := 0

	p.DstBroadcast = buf[i] == 1
	i++
	p.DstAddrType = int(buf[i])
	i++

	// Next six bytes encode either a mac address or an ip address and port
	switch p.DstAddrType {
	case addrTypeEth:
		p.DstMac = buf[i : i+6]

	case addrTypeUDP:
		p.DstIP = buf[i : i+4]
		port := binary.BigEndian.Uint16(buf[i+4 : i+6])
		p.DstPort = uint(port)

	default:
		return nil, errors.New("Unknown destination address type")
	}
	i += 6

	p.SrcBroadcast = buf[i] == 1
	i++
	p.SrcAddrType = int(buf[i])
	i++

	// Next six bytes encode either a mac address or an ip address and port
	switch p.SrcAddrType {
	case addrTypeEth:
		p.SrcMac = buf[i : i+6]

	case addrTypeUDP:
		p.SrcIP = buf[i : i+4]
		port := binary.BigEndian.Uint16(buf[i+4 : i+6])
		p.SrcPort = uint(port)

	default:
		return nil, errors.New("Unknown source address type")
	}
	i += 6

	p.Seq = int(binary.BigEndian.Uint16(buf[i : i+2]))
	i += 2

	p.UdapType = int(binary.BigEndian.Uint16(buf[i : i+2]))
	i += 2

	p.UcpFlags = buf[i]
	i++

	p.UapClass = buf[i : i+4]
	i += 4

	p.UcpMethod = int(binary.BigEndian.Uint16(buf[i : i+2]))
	i += 2

	// Remaining data is returned as-is
	p.Data = buf[i:]

	return &p, nil
}

func (p packet) assemble() []byte {
	var buf []byte
	portSlice := make([]byte, 2)

	// Destination
	buf = append(buf, func() byte {
		if p.DstBroadcast {
			return 0x01
		}
		return 0x00
	}())
	buf = append(buf, byte(p.DstAddrType))
	switch p.DstAddrType {
	case addrTypeEth:
		buf = append(buf, p.DstMac...)
	case addrTypeUDP:
		buf = append(buf, p.DstIP...)
		binary.BigEndian.PutUint16(portSlice, uint16(p.DstPort))
		buf = append(buf, portSlice...)
	}

	// Source
	buf = append(buf, func() byte {
		if p.SrcBroadcast {
			return 0x01
		}
		return 0x00
	}())
	buf = append(buf, byte(p.SrcAddrType))
	switch p.SrcAddrType {
	case addrTypeEth:
		buf = append(buf, p.SrcMac...)
	case addrTypeUDP:
		buf = append(buf, p.SrcIP...)
		binary.BigEndian.PutUint16(portSlice, uint16(p.SrcPort))
		buf = append(buf, portSlice...)
	}

	buf = append(buf, []byte{0x00, 0x01}...)
	buf = append(buf, udapTypeUCP...)
	buf = append(buf, []byte{0x01}...)
	buf = append(buf, uapClassUCP...)
	buf = append(buf, 0x00, byte(p.UcpMethod))

	switch p.UcpMethod {
	case ucpMethodGetData:
		buf = append(buf, defaultCredentials...)
		buf = append(buf, 0x00, byte(len(dataItems)))
		for i, v := range dataItems {
			offset := make([]byte, 2)
			length := make([]byte, 2)
			binary.BigEndian.PutUint16(offset, uint16(i))
			binary.BigEndian.PutUint16(length, uint16(v.length))
			buf = append(buf, offset...)
			buf = append(buf, length...)
		}
		buf = append(buf, 0x00)

	case ucpMethodSetData:
		buf = append(buf, defaultCredentials...)
		buf = append(buf, p.Data...)
	}

	return buf
}

type fields map[byte][]byte

func (p packet) parseFields() (fields, error) {
	// Data format is a repeated list of:
	//  UCP Code (1 byte)
	//  Length (1 byte)
	//  Data []byte
	//  If Length is zero, it is the end of the data

	data := make(fields)
	buf := p.Data
	for {
		if len(buf) < 2 {
			break
		}
		ucpCode := buf[0]
		length := int(buf[1])
		if length == 0 {
			break
		}
		if len(buf) < length+2 {
			return nil, errors.New("Data buffer too short")
		}
		data[ucpCode] = buf[2 : length+2]
		buf = buf[length+2:]
	}

	return data, nil
}

func (p packet) parseData() (map[int][]byte, error) {
	data := make(map[int][]byte)
	buf := p.Data
	if len(buf) < 2 {
		return nil, errors.New("No data")
	}
	numValues := int(binary.BigEndian.Uint16(buf[0:2]))
	buf = buf[2:]
	for i := 0; i < numValues; i++ {
		if len(buf) < 4 {
			break
		}
		index := int(binary.BigEndian.Uint16(buf[0:2]))
		length := int(binary.BigEndian.Uint16(buf[2:4]))
		if len(buf) < 4+length {
			return nil, errors.New("Data format error")
		}
		data[index] = buf[4 : 4+length]
		buf = buf[4+length:]
	}

	return data, nil
}

func setDeviceData(sb *Sb, data fields) {
	for i, v := range data {
		switch i {
		//case UCPCodeZero:
		//case UCPCodeOne:
		case UCPCodeDeviceName:
			sb.Name = string(v)
		case UCPCodeDeviceType:
			sb.Type = string(v)
		// case UCPodeUseDHCP    :
		case UCPCodeIPAddr:
			sb.IPAddr = v
		case UCPCodeSubnetMask:
			sb.SubnetMask = v
		case UCPCodeGatewayAddr:
			sb.GatewayAddr = v
		// case UCPCodeEight       :
		case UCPCodeFirmwareRev:
			sb.FirmwareRev = uint(binary.BigEndian.Uint16(v))
		case UCPCodeHardwareRev:
			sb.HardwareRev = uint(binary.BigEndian.Uint32(v))
		case UCPCodeDeviceID:
			sb.ID = uint(binary.BigEndian.Uint16(v))
		case UCPCodeDeviceStatus:
			sb.Status = string(v)
			// case UCPCodeUUID :
		}
	}
}
