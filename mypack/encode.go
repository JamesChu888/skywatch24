package mypack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func encode(input interface{}) ([]byte, error) {
	var buf bytes.Buffer

	switch val := input.(type) {

	// nil
	case nil:
		buf.WriteByte(leadingByteNil)

	// int
	case int, int8, int16, int32, int64:
		encodeInt(val, &buf)

	// uint
	case uint, uint8, uint16, uint32, uint64:
		encodeUint(val, &buf)

	// string
	case string:
		encodeString(val, &buf)

	// bool
	case bool:
		encodeBool(val, &buf)

	// float
	case float32, float64:
		encodeFloat(val, &buf)

	// array
	case []interface{}:
		encodeArray(val, &buf)

	// map
	case map[string]interface{}:
		encodeMap(val, &buf)

	// binary
	case []byte:
		encodeBinary(val, &buf)

	default:
		return nil, fmt.Errorf("encode error msg: %v", reflect.TypeOf(input))
	}

	return buf.Bytes(), nil
}

// int
func encodeInt(input interface{}, buf *bytes.Buffer) {
	switch val := input.(type) {

	case int8:
		buf.WriteByte(leadingByteInt8)
		buf.WriteByte(byte(val))

	case int16:
		buf.WriteByte(leadingByteInt16)
		binary.Write(buf, binary.BigEndian, val)

	case int32:
		buf.WriteByte(leadingByteInt32)
		binary.Write(buf, binary.BigEndian, val)

	case int64:
		buf.WriteByte(leadingByteInt64)
		binary.Write(buf, binary.BigEndian, val)

	case int:
		if val >= -32 && val <= 127 { // fixint
			buf.WriteByte(byte(val))
		} else if val >= -128 && val < -32 {
			buf.WriteByte(leadingByteInt8)
			buf.WriteByte(byte(val))
		} else if val >= -32768 && val < -128 {
			buf.WriteByte(leadingByteInt16)
			binary.Write(buf, binary.BigEndian, int16(val))
		} else if val >= -2147483648 && val < -32768 {
			buf.WriteByte(leadingByteInt32)
			binary.Write(buf, binary.BigEndian, int32(val))
		} else {
			buf.WriteByte(leadingByteInt64)
			binary.Write(buf, binary.BigEndian, int64(val))
		}
	}
}

// uint
func encodeUint(input interface{}, buf *bytes.Buffer) {
	switch val := input.(type) {
	case uint:
		if val <= 127 { // fixuint
			buf.WriteByte(byte(val))
		} else {
			buf.WriteByte(leadingByteUint64)
			binary.Write(buf, binary.BigEndian, uint64(val))
		}
	case uint8:
		buf.WriteByte(leadingByteUint8)
		buf.WriteByte(byte(val))
	case uint16:
		buf.WriteByte(leadingByteUint16)
		binary.Write(buf, binary.BigEndian, val)
	case uint32:
		buf.WriteByte(leadingByteUint32)
		binary.Write(buf, binary.BigEndian, val)
	case uint64:
		buf.WriteByte(leadingByteUint64)
		binary.Write(buf, binary.BigEndian, val)
	}
}

// string
func encodeString(input string, buf *bytes.Buffer) {
	length := len(input)

	if length <= fixStrMaxLength {
		buf.WriteByte(fixStrStart | byte(length))
	} else if length <= str8MaxLength {
		buf.WriteByte(leadingByteStr8)
		buf.WriteByte(byte(length))
	} else if length <= str16MaxLength {
		buf.WriteByte(leadingByteStr16)
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else {
		buf.WriteByte(leadingByteStr32)
		binary.Write(buf, binary.BigEndian, uint32(length))
	}

	buf.WriteString(input)
}

// bool
func encodeBool(input bool, buf *bytes.Buffer) {
	if input {
		buf.WriteByte(leadingByteTrue)
	} else {
		buf.WriteByte(leadingByteFalse)
	}
}

// float
func encodeFloat(input interface{}, buf *bytes.Buffer) {
	switch val := input.(type) {
	case float32:
		buf.WriteByte(leadingByteFloat32)
		binary.Write(buf, binary.BigEndian, val)
	case float64:
		buf.WriteByte(leadingByteFloat64)
		binary.Write(buf, binary.BigEndian, val)
	}
}

// array (may have nested element, need recursive process)
func encodeArray(input []interface{}, buf *bytes.Buffer) {
	length := len(input)
	if length <= fixArrayLength {
		buf.WriteByte(fixArrayStart | byte(length))
	} else if length <= array16MaxLength {
		buf.WriteByte(leadingByteArray16)
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else {
		buf.WriteByte(leadingByteArray32)
		binary.Write(buf, binary.BigEndian, uint32(length))
	}
	for _, item := range input {
		itemEncoded, err := encode(item)
		if err != nil {
			panic(err)
		}
		buf.Write(itemEncoded)
	}
}

// map (may have nested element, need recursive process)
func encodeMap(input map[string]interface{}, buf *bytes.Buffer) {
	length := len(input)
	if length <= fixMapMaxLength {
		buf.WriteByte(fixMapStart | byte(length))
	} else if length <= map16MaxLength {
		buf.WriteByte(leadingByteMap16)
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else {
		buf.WriteByte(leadingByteMap32)
		binary.Write(buf, binary.BigEndian, uint32(length))
	}
	for key, value := range input {
		keyEncoded, err := encode(key)
		if err != nil {
			panic(err)
		}
		buf.Write(keyEncoded)

		valueEncoded, err := encode(value)
		if err != nil {
			panic(err)
		}
		buf.Write(valueEncoded)
	}
}

// binary
func encodeBinary(input []byte, buf *bytes.Buffer) {
	length := len(input)
	if length <= bin8MaxLength {
		buf.WriteByte(leadingByteBin8)
		buf.WriteByte(byte(length))
	} else if length <= bin16MaxLength {
		buf.WriteByte(leadingByteBin16)
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else {
		buf.WriteByte(leadingByteBin32)
		binary.Write(buf, binary.BigEndian, uint32(length))
	}
	buf.Write(input)
}
