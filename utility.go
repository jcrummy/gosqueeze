package gosqueeze

import (
	"errors"
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
