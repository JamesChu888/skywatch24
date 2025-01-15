package mypack

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMypack(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		// FixInt (Positive and Negative)
		{"+fixint", int8(42), int8(42)},
		{"-fixint", int8(-31), int8(-31)},

		// int
		{"int8", int8(-128), int8(-128)},
		{"int16", int16(-32768), int16(-32768)},
		{"int32", int32(-2147483648), int32(-2147483648)},
		{"int64", int64(-9223372036854775808), int64(-9223372036854775808)},

		// uint
		{"uint8", uint8(255), uint8(255)},
		{"uint16", uint16(65535), uint16(65535)},
		{"uint32", uint32(4294967295), uint32(4294967295)},
		{"uint64", uint64(18446744073709551615), uint64(18446744073709551615)},

		// float
		{"float32", float32(3.14), float32(3.14)},
		{"float64", float64(2.71828), float64(2.71828)},

		// bool
		{"bool_true", true, true},
		{"bool_false", false, false},

		// string
		{"string", "hello world", "hello world"},

		// binary
		{"binary", []byte{0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03}},

		// array
		{"array", []interface{}{1, "two", true}, []interface{}{int8(1), "two", true}},

		// map
		{"map", map[string]interface{}{"key": "value", "num": 42}, map[string]interface{}{"key": "value", "num": int8(42)}},

		// nil
		{"nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 編碼
			encoded, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// 解碼
			buf := bytes.NewReader(encoded)
			decoded, err := Decode(buf)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// 驗證結果
			if !reflect.DeepEqual(decoded, tt.expected) {
				t.Errorf("%v (%T) <> %v (%T)", tt.expected, tt.expected, decoded, decoded)
			}
		})
	}
}
