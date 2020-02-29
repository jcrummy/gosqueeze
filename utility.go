package gosqueeze

import (
	"errors"
	"net"
	"reflect"
	"strconv"
	"strings"
)

func getTag(field reflect.StructField) (int, int, error) {
	itemTag := field.Tag.Get("gosqueeze")
	if itemTag == "" {
		return 0, 0, errors.New("No tag attached to struct field")
	}
	tagValues := strings.Split(itemTag, ",")
	if len(tagValues) < 2 {
		return 0, 0, errors.New("Field tag requires two values (offset,length)")
	}
	offset, err := strconv.Atoi(tagValues[0])
	if err != nil {
		return 0, 0, errors.New("Offset is not a number: " + err.Error())
	}
	length, err := strconv.Atoi(tagValues[1])
	if err != nil {
		return 0, 0, errors.New("Length is not a number: " + err.Error())
	}
	return offset, length, nil
}

func getFieldOffsets(field reflect.StructField) (int, string, error) {
	itemTag := field.Tag.Get("gosqueeze")
	if itemTag == "" {
		return 0, "", errors.New("No tag attached to struct field")
	}
	tagValues := strings.Split(itemTag, ",")
	if len(tagValues) < 2 {
		return 0, "", errors.New("Field tag requires two values (offset,length)")
	}
	offset, err := strconv.Atoi(tagValues[0])
	if err != nil {
		return 0, "", errors.New("Offset is not a number: " + err.Error())
	}
	return offset, field.Name, nil
}

func getOffsetMap(dataFields interface{}) map[int]string {
	ret := make(map[int]string)
	st := reflect.TypeOf(dataFields)
	for i := 0; i < st.Elem().NumField(); i++ {
		offset, name, err := getFieldOffsets(st.Elem().Field(i))
		if err != nil {
			continue
		}
		ret[offset] = name
	}
	return ret
}

func pack(v interface{}, length int) []byte {
	s := reflect.TypeOf(v)
	slice := make([]byte, length)
	switch s.String() {
	case "bool":
		if v.(bool) {
			slice[0] = 1
			break
		}
		slice[0] = 0

	case "string":
		copy(slice, []byte(v.(string)))

	case "uint8":
		copy(slice, []byte{v.(uint8)})

	case "[]uint8":
		copy(slice, []byte(v.([]uint8)))

	case "net.IP":
		copy(slice, v.(net.IP))
	}
	if len(slice) < length {
		slice = append(slice, make([]byte, 33-len(slice))...)
	} else if len(slice) > length {
		slice = slice[:length]
	}
	return slice
}
