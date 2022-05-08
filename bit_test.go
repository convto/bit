package bit

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type encDecTest struct {
	enc string
	dec []byte
}

var encDecTests = []encDecTest{
	{"", []byte{}},
	{"0000000000000001000000100000001100000100000001010000011000000111", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"0000100000001001000010100000101100001100000011010000111000001111", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"1111000011110001111100101111001111110100111101011111011011110111", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"1111100011111001111110101111101111111100111111011111111011111111", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"01100111", []byte{'g'}},
	{"1110001110100001", []byte{0b11100011, 0b10100001}},
}

func TestEncode(t *testing.T) {
	for i, test := range encDecTests {
		dst := make([]byte, EncodedLen(len(test.dec)))
		n := Encode(dst, test.dec)
		if n != len(dst) {
			t.Errorf("#%d: bad return value: got: %d want: %d", i, n, len(dst))
		}
		if string(dst) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.enc)
		}
	}
}

func TestDecode(t *testing.T) {
	for i, test := range encDecTests {
		dst := make([]byte, DecodedLen(len(test.enc)))
		n, err := Decode(dst, []byte(test.enc))
		if err != nil {
			t.Errorf("#%d: bad return value: got:%d want:%d", i, n, len(dst))
		} else if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.dec)
		}
	}
}

func TestEncodeToString(t *testing.T) {
	for i, test := range encDecTests {
		s := EncodeToString(test.dec)
		if s != test.enc {
			t.Errorf("#%d got:%s want:%s", i, s, test.enc)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for i, test := range encDecTests {
		dst, err := DecodeString(test.enc)
		if err != nil {
			t.Errorf("#%d: unexpected err value: %s", i, err)
			continue
		}
		if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: #%v", i, dst, test.dec)
		}
	}
}

var errTests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
	{"1111111", "", ErrLength},
	{"z1111111", "", InvalidByteError('z')},
	{"11111111z", "\xff", InvalidByteError('z')},
	{"111111110", "\xff", ErrLength},
	{"1111111g", "", InvalidByteError('g')},
	{"11111111gg", "\xff", InvalidByteError('g')},
	{"0000000\x01", "", InvalidByteError('\x01')},
	{"11111111000000001111111", "\xff\x00", ErrLength},
}

func TestDecodeErr(t *testing.T) {
	for _, tt := range errTests {
		out := make([]byte, len(tt.in)+10)
		n, err := Decode(out, []byte(tt.in))
		if string(out[:n]) != tt.out || err != tt.err {
			t.Errorf("Decode(%q) = %q, %v, want %q, %v", tt.in, string(out[:n]), err, tt.out, tt.err)
		}
	}
}

func TestDecodeStringErr(t *testing.T) {
	for _, tt := range errTests {
		out, err := DecodeString(tt.in)
		if string(out) != tt.out || err != tt.err {
			t.Errorf("DecodeString(%q) = %q, %v, want %q, %v", tt.in, out, err, tt.out, tt.err)
		}
	}
}

func TestEncoderDecoder(t *testing.T) {
	for _, multiplier := range []int{1, 128, 192} {
		for _, test := range encDecTests {
			input := bytes.Repeat(test.dec, multiplier)
			output := strings.Repeat(test.enc, multiplier)

			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			r := struct{ io.Reader }{bytes.NewReader(input)} // io.Reader only; not io.WriterTo
			if n, err := io.CopyBuffer(enc, r, make([]byte, 7)); n != int64(len(input)) || err != nil {
				t.Errorf("encoder.Write(%q*%d) = (%d, %v), want (%d, nil)", test.dec, multiplier, n, err, len(input))
				continue
			}

			if encDst := buf.String(); encDst != output {
				t.Errorf("buf(%q*%d) = %v, want %v", test.dec, multiplier, encDst, output)
				continue
			}

			dec := NewDecoder(&buf)
			var decBuf bytes.Buffer
			w := struct{ io.Writer }{&decBuf} // io.Writer only; not io.ReaderFrom
			if _, err := io.CopyBuffer(w, dec, make([]byte, 7)); err != nil || decBuf.Len() != len(input) {
				t.Errorf("decoder.Read(%q*%d) = (%d, %v), want (%d, nil)", test.enc, multiplier, decBuf.Len(), err, len(input))
			}

			if !bytes.Equal(decBuf.Bytes(), input) {
				t.Errorf("decBuf(%q*%d) = %v, want %v", test.dec, multiplier, decBuf.Bytes(), input)
				continue
			}
		}
	}
}

func TestDecoderErr(t *testing.T) {
	for _, tt := range errTests {
		dec := NewDecoder(strings.NewReader(tt.in))
		out, err := io.ReadAll(dec)
		wantErr := tt.err
		// Decoder is reading from stream, so it reports io.ErrUnexpectedEOF instead of ErrLength.
		if wantErr == ErrLength {
			wantErr = io.ErrUnexpectedEOF
		}
		if string(out) != tt.out || err != wantErr {
			t.Errorf("NewDecoder(%q) = %q, %v, want %q, %v", tt.in, out, err, tt.out, wantErr)
		}
	}
}

func TestDumper(t *testing.T) {
	var in [40]byte
	for i := range in {
		in[i] = byte(i + 30)
	}

	for stride := 1; stride < len(in); stride++ {
		var out bytes.Buffer
		dumper := Dumper(&out)
		done := 0
		for done < len(in) {
			todo := done + stride
			if todo > len(in) {
				todo = len(in)
			}
			dumper.Write(in[done:todo])
			done = todo
		}

		dumper.Close()
		if !bytes.Equal(out.Bytes(), expectedBinDump) {
			t.Errorf("stride: %d failed. got:\n%s\nwant:\n%s", stride, out.Bytes(), expectedBinDump)
		}
	}
}

func TestDumper_doubleclose(t *testing.T) {
	var out bytes.Buffer
	dumper := Dumper(&out)

	dumper.Write([]byte(`gopher`))
	dumper.Close()
	dumper.Close()
	dumper.Write([]byte(`gopher`))
	dumper.Close()

	expected := "00000000: 01100111 01101111 01110000 01101000 01100101 01110010  gopher\n"
	if out.String() != expected {
		t.Fatalf("got:\n%#v\nwant:\n%#v", out.String(), expected)
	}
}

func TestDumper_earlyclose(t *testing.T) {
	var out bytes.Buffer
	dumper := Dumper(&out)

	dumper.Close()
	dumper.Write([]byte(`gopher`))

	expected := ""
	if out.String() != expected {
		t.Fatalf("got:\n%#v\nwant:\n%#v", out.String(), expected)
	}
}

func TestDump(t *testing.T) {
	var in [40]byte
	for i := range in {
		in[i] = byte(i + 30)
	}

	out := []byte(Dump(in[:]))
	if !bytes.Equal(out, expectedBinDump) {
		t.Errorf("got:\n%s\nwant:\n%s", out, expectedBinDump)
	}
}

var expectedBinDump = []byte(`00000000: 00011110 00011111 00100000 00100001 00100010 00100011  .. !"#
00000006: 00100100 00100101 00100110 00100111 00101000 00101001  $%&'()
0000000c: 00101010 00101011 00101100 00101101 00101110 00101111  *+,-./
00000012: 00110000 00110001 00110010 00110011 00110100 00110101  012345
00000018: 00110110 00110111 00111000 00111001 00111010 00111011  6789:;
0000001e: 00111100 00111101 00111110 00111111 01000000 01000001  <=>?@A
00000024: 01000010 01000011 01000100 01000101                    BCDE
`)
