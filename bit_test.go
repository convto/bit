package bit

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
	}
	type want struct {
		n       int
		encoded []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test bin",
			args: args{
				dst: make([]byte, EncodedLen(9)),
				src: []byte{0x08, 0xb9, 0x60, 0x10, 0xb2, 0x92, 0x04, 0x18, 0x01},
			},
			want: want{
				n:       EncodedLen(9),
				encoded: []byte("000010001011100101100000000100001011001010010010000001000001100000000001"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := Encode(tt.args.dst, tt.args.src); gotN != tt.want.n {
				t.Errorf("Encode() = %v, want.n %v", gotN, tt.want.n)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want.encoded) {
				t.Errorf("Encode() = %v, want.encoded %v", tt.args.dst, tt.want.encoded)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
	}
	type want struct {
		n       int
		err     bool
		decoded []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test bin",
			args: args{
				dst: make([]byte, DecodedLen(72)),
				src: []byte("000010001011100101100000000100001011001010010010000001000001100000000001"),
			},
			want: want{
				n:       DecodedLen(72),
				decoded: []byte{0x08, 0xb9, 0x60, 0x10, 0xb2, 0x92, 0x04, 0x18, 0x01},
			},
		},
		{
			name: "invalid byte",
			args: args{
				dst: []byte{},
				src: []byte("Z"),
			},
			want: want{
				n:       0,
				decoded: []byte{},
				err:     true,
			},
		},
		{
			name: "invalid length",
			args: args{
				dst: make([]byte, DecodedLen(15)),
				src: []byte("010101010101010"),
			},
			want: want{
				n:       1,
				decoded: []byte{0x55},
				err:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := Decode(tt.args.dst, tt.args.src)
			if (err != nil) != tt.want.err {
				t.Errorf("Decode() error = %v, want.err %v", err, tt.want.err)
				return
			}
			if gotN != tt.want.n {
				t.Errorf("Decode() gotN = %v, want.n %v", gotN, tt.want.n)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want.decoded) {
				t.Errorf("Decode() = %v, want.decoded %v", tt.args.dst, tt.want.decoded)
			}
		})
	}
}
