package mypack

import "bytes"

const (
	leadingByteInt8  = 0xd0
	leadingByteInt16 = 0xd1
	leadingByteInt32 = 0xd2
	leadingByteInt64 = 0xd3

	leadingByteUint8  = 0xcc
	leadingByteUint16 = 0xcd
	leadingByteUint32 = 0xce
	leadingByteUint64 = 0xcf

	leadingByteFloat32 = 0xca
	leadingByteFloat64 = 0xcb

	fixStrMaxLength = 31
	fixStrStart     = 0xa0 // 0xa0 - 0xbf
	fixStrEnd       = 0xbf

	str8MaxLength  = 255
	str16MaxLength = 65535

	leadingByteStr8  = 0xd9
	leadingByteStr16 = 0xda
	leadingByteStr32 = 0xdb

	leadingByteFalse = 0xc2
	leadingByteTrue  = 0xc3

	leadingByteNil = 0xc0

	fixArrayLength = 15
	fixArrayStart  = 0x90 // 0x90 ~ 0x9f
	fixArrayEnd    = 0x9f

	array16MaxLength = 65535

	leadingByteArray16 = 0xdc
	leadingByteArray32 = 0xdd

	fixMapMaxLength = 15
	fixMapStart     = 0x80 // 0x80 ~ 0x8f
	fixMapEnd       = 0x8f

	map16MaxLength = 65535

	leadingByteMap16 = 0xde
	leadingByteMap32 = 0xdf

	bin8MaxLength  = 255
	bin16MaxLength = 65535

	leadingByteBin8  = 0xc4
	leadingByteBin16 = 0xc5
	leadingByteBin32 = 0xc6

	leadingByteFixExt1  = 0xd4
	leadingByteFixExt2  = 0xd5
	leadingByteFixExt4  = 0xd6
	leadingByteFixExt8  = 0xd7
	leadingByteFixExt16 = 0xd8
	leadingByteExt8     = 0xc7
	leadingByteExt16    = 0xc8
	leadingByteExt32    = 0xc9
)

func Marshal(input interface{}) ([]byte, error) {
	return encode(input)
}

func Unmarshal(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	return decode(buf)
}
