package bit

import "io"

const (
	bitTable = "01"
	b1 = 0b00000001
	b2 = 0b00000010
	b3 = 0b00000100
	b4 = 0b00001000
	b5 = 0b00010000
	b6 = 0b00100000
	b7 = 0b01000000
)

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 8.
func EncodedLen(n int) int {return n*8}

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements bits encoding.
func Encode(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = bitTable[v>>7]
		dst[j+1] = bitTable[(v&b7)>>6]
		dst[j+2] = bitTable[(v&b6)>>5]
		dst[j+3] = bitTable[(v&b5)>>4]
		dst[j+4] = bitTable[(v&b4)>>3]
		dst[j+5] = bitTable[(v&b3)>>2]
		dst[j+6] = bitTable[(v&b2)>>1]
		dst[j+7] = bitTable[(v&b1)]
		j += 8
	}
	return len(src)*8
}

// EncodeToString returns the bits encoding of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

// bufferSize is the number of bits characters to buffer in encoder.
const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns an io.Writer that writes bits characters to w.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / 2
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var written int
		encoded := Encode(e.out[:], p[:chunkSize])
		written, e.err = e.w.Write(e.out[:encoded])
		n += written / 8
		p = p[chunkSize:]
	}
	return n, e.err
}
