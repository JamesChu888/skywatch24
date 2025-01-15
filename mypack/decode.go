package mypack

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Decode(buf *bytes.Reader) (interface{}, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {

	// nil
	case b == leadingByteNil:
		return nil, nil

	// uint
	case b >= leadingByteUint8 && b <= leadingByteUint64:
		return decodeUint(b, buf)

	// positive fixint: 0x00 ~ 0x7f / 0 ~ 127
	// negative fixint: 0xe0 ~ 0xff / -32 ~ -1
	// int
	case /* b >= 0x00 && */ b <= 0x7f, b >= 0xe0 /* && b <= 0xff */, b >= leadingByteInt8 && b <= leadingByteInt64:
		return decodeInt(b, buf)

	// float
	case b == leadingByteFloat32 || b == leadingByteFloat64:
		return decodeFloat(b, buf)

	// false
	case b == leadingByteFalse:
		return false, nil

	// true
	case b == leadingByteTrue:
		return true, nil

	// string
	case b >= fixStrStart && b <= fixStrEnd, b == leadingByteStr8 || b == leadingByteStr16 || b == leadingByteStr32:
		return decodeString(b, buf)

	// array
	case b >= fixArrayStart && b <= fixArrayEnd, b == leadingByteArray16 || b == leadingByteArray32:
		return decodeArray(b, buf)

	// map
	case b >= fixMapStart && b <= fixMapEnd, b == leadingByteMap16 || b == leadingByteMap32:
		return decodeMap(b, buf)

	// Binary
	case b >= leadingByteBin8 && b <= leadingByteBin32:
		return decodeBinary(b, buf)

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}
}

// int
func decodeInt(b byte, buf *bytes.Reader) (interface{}, error) {
	switch {
	// fixint
	case /* b >= 0x00 && */ b <= 0x7f:
		return int8(b), nil

	// negative fixint
	case b >= 0xe0 /* && b <= 0xff */ :
		return int8(int16(b) - int16(0x100)), nil

	// int8
	case b == leadingByteInt8:
		var value int8
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	// int16
	case b == leadingByteInt16:
		var value int16
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	// int32
	case b == leadingByteInt32:
		var value int32
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	// int64
	case b == leadingByteInt64:
		var value int64
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}
}

// uint
func decodeUint(b byte, buf *bytes.Reader) (interface{}, error) {
	switch b {
	case leadingByteUint8:
		var value uint8
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return uint8(value), nil
	case leadingByteUint16:
		var value uint16
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return uint16(value), nil
	case leadingByteUint32:
		var value uint32
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return uint32(value), nil
	case leadingByteUint64:
		var value uint64
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return uint64(value), nil
	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}
}

// float
func decodeFloat(b byte, buf *bytes.Reader) (interface{}, error) {
	switch b {
	case leadingByteFloat32:
		var value float32
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	case leadingByteFloat64:
		var value float64
		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return nil, err
		}
		return value, nil

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}
}

// string
func decodeString(b byte, buf *bytes.Reader) (interface{}, error) {
	var length int
	switch {
	// fixstr: 長度 0 到 31
	case b >= fixStrStart && b <= fixStrEnd:
		length = int(b - fixStrStart)

	// str8
	case b == leadingByteStr8:
		var l uint8
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	// str16
	case b == leadingByteStr16:
		var l uint16
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	// str32:
	case b == leadingByteStr32:
		var l uint32
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}

	data := make([]byte, length)
	if _, err := buf.Read(data); err != nil {
		return nil, err
	}
	return string(data), nil
}

// array (need recursive decode)
func decodeArray(b byte, buf *bytes.Reader) (interface{}, error) {
	var length int

	switch {
	// fixarray: 長度 0 到 15
	case b >= fixArrayStart && b <= fixArrayEnd:
		length = int(b - fixArrayStart)

	// array16
	case b == leadingByteArray16:
		var l uint16
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	// array32
	case b == leadingByteArray32:
		var l uint32
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}

	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		item, err := Decode(buf)
		if err != nil {
			return nil, err
		}
		array[i] = item

		fmt.Printf("item: %v\n", item)
	}
	return array, nil
}

// map (need recursive decode)
func decodeMap(b byte, buf *bytes.Reader) (interface{}, error) {
	var length int

	switch {
	// fixmap: 長度 0 到 15
	case b >= fixMapStart && b <= fixMapEnd:
		length = int(b - fixMapStart)

	// map16
	case b == leadingByteMap16:
		var l uint16
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	// map32
	case b == leadingByteMap32:
		var l uint32
		if err := binary.Read(buf, binary.BigEndian, &l); err != nil {
			return nil, err
		}
		length = int(l)

	default:
		return nil, fmt.Errorf("type undefined: %v", b)
	}

	mapping := make(map[string]interface{}, length)
	for i := 0; i < length; i++ {
		key, err := Decode(buf)
		if err != nil {
			return nil, err
		}
		value, err := Decode(buf)
		if err != nil {
			return nil, err
		}
		mapping[key.(string)] = value
	}
	return mapping, nil
}

// binary
func decodeBinary(b byte, buf *bytes.Reader) ([]byte, error) {
	var length int
	switch b {
	// bin8
	case leadingByteBin8:
		var tmp uint8
		if err := binary.Read(buf, binary.BigEndian, &tmp); err != nil {
			return nil, err
		}
		length = int(tmp)

	// bin16
	case leadingByteBin16:
		var tmp uint16
		if err := binary.Read(buf, binary.BigEndian, &tmp); err != nil {
			return nil, err
		}
		length = int(tmp)

	// bin32
	case leadingByteBin32:
		var tmp uint32
		if err := binary.Read(buf, binary.BigEndian, &tmp); err != nil {
			return nil, err
		}
		length = int(tmp)
	}

	data := make([]byte, length)
	if _, err := buf.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
