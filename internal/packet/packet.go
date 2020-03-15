// Copyright 2020 John Crummy. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package packet

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"reflect"

	"github.com/jcrummy/gosqueeze/internal/constants"
	"github.com/jcrummy/gosqueeze/internal/util"
)

// Packet represents a SqueezeBox configuration packet. The same
// format is used for requests and replies.
type Packet struct {
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

// Parse returns a packet struct from the raw byte slice.
func Parse(buf []byte) (*Packet, error) {
	if len(buf) < 27 {
		return nil, errors.New("Packet length too short")
	}

	var p Packet
	i := 0

	p.DstBroadcast = buf[i] == 1
	i++
	p.DstAddrType = int(buf[i])
	i++

	// Next six bytes encode either a mac address or an ip address and port
	switch p.DstAddrType {
	case constants.AddrTypeEth:
		p.DstMac = buf[i : i+6]

	case constants.AddrTypeUDP:
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
	case constants.AddrTypeEth:
		p.SrcMac = buf[i : i+6]

	case constants.AddrTypeUDP:
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

// Assemble provides a raw byte slice ready to send over the network.
func (p Packet) Assemble() []byte {
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
	case constants.AddrTypeEth:
		buf = append(buf, p.DstMac...)
	case constants.AddrTypeUDP:
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
	case constants.AddrTypeEth:
		buf = append(buf, p.SrcMac...)
	case constants.AddrTypeUDP:
		buf = append(buf, p.SrcIP...)
		binary.BigEndian.PutUint16(portSlice, uint16(p.SrcPort))
		buf = append(buf, portSlice...)
	}

	buf = append(buf, []byte{0x00, 0x01}...)
	buf = append(buf, constants.UdapTypeUCP...)
	buf = append(buf, []byte{0x01}...)
	buf = append(buf, constants.UapClassUCP...)
	buf = append(buf, 0x00, byte(p.UcpMethod))

	switch p.UcpMethod {
	case constants.UCPMethodGetData:
		buf = append(buf, constants.DefaultCredentials...)
		buf = append(buf, p.Data...)

	case constants.UCPMethodSetData:
		buf = append(buf, constants.DefaultCredentials...)
		buf = append(buf, p.Data...)
	}

	return buf
}

// Fields is a map of configuration data points in their raw format.
type Fields map[byte][]byte

// ParseFields returns a field map of raw field data from the .Data
// byte slice of the packet.
func (p Packet) ParseFields() (Fields, error) {
	// Data format is a repeated list of:
	//  UCP Code (1 byte)
	//  Length (1 byte)
	//  Data []byte
	//  If Length is zero, it is the end of the data

	data := make(Fields)
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

// ParseData populates a struct based on the .Data byte slice
// of the packet. Field data is entered based on the tagged offset
// value of the structure.
func (p Packet) ParseData(dataFields interface{}) error {
	fieldOffsets := util.GetOffsetMap(dataFields)
	s := reflect.ValueOf(dataFields).Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("Not a structure")
	}

	buf := p.Data
	if len(buf) < 2 {
		return errors.New("No data")
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
			return errors.New("Data format error")
		}
		name, ok := fieldOffsets[index]
		if !ok {
			log.Printf("Field not found for offset %d.\n", index)
			continue
		}
		f := s.FieldByName(name)
		if !f.IsValid() {
			log.Printf("Field not found for offset %d: %s.\n", index, name)
			continue
		}
		if !f.CanSet() {
			log.Printf("Can't set this data for %+v\n", name)
			buf = buf[4+length:]
			continue
		}
		data := buf[4 : 4+length]

		switch f.Type().String() {
		case "bool":
			f.SetBool(data[0] == 0x01)

		case "string":
			f.SetString(string(data))

		case "uint8":
			f.SetUint(uint64(data[0]))

		case "net.IP":
			f.SetBytes(data)
		}
		buf = buf[4+length:]
	}

	return nil
}

// SetDataRetrieve applies the set of configuration values to be retrieved
func (p *Packet) SetDataRetrieve(dataFields interface{}) error {
	st := reflect.TypeOf(dataFields)
	p.Data = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Data, uint16(st.NumField()))
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		offset, length, err := util.GetTag(field)
		if err != nil {
			continue
		}
		offsetB := make([]byte, 2)
		lengthB := make([]byte, 2)
		binary.BigEndian.PutUint16(offsetB, uint16(offset))
		binary.BigEndian.PutUint16(lengthB, uint16(length))
		p.Data = append(p.Data, offsetB...)
		p.Data = append(p.Data, lengthB...)
	}
	return nil
}

// SetDataForSave applies configuration data to save to the device
func (p *Packet) SetDataForSave(dataFields interface{}) int {
	st := reflect.TypeOf(dataFields)
	p.Data = make([]byte, 2)
	binary.BigEndian.PutUint16(p.Data, uint16(st.NumField()))
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		offset, length, err := util.GetTag(field)
		if err != nil {
			continue
		}
		offsetB := make([]byte, 2)
		lengthB := make([]byte, 2)
		binary.BigEndian.PutUint16(offsetB, uint16(offset))
		binary.BigEndian.PutUint16(lengthB, uint16(length))
		p.Data = append(p.Data, offsetB...)
		p.Data = append(p.Data, lengthB...)
		fieldValue := reflect.ValueOf(dataFields).FieldByName(field.Name)
		p.Data = append(p.Data, util.Pack(fieldValue.Interface(), length)...)
	}
	return st.NumField()
}
