package bit

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	bitTable = "01"
	b1       = 0b00000001
	b2       = 0b00000010
	b3       = 0b00000100
	b4       = 0b00001000
	b5       = 0b00010000
	b6       = 0b00100000
	b7       = 0b01000000
)

// ErrLength reports an attempt to decode an odd-length input
// using Decode or DecodeString.
// The stream-based Decoder returns io.ErrUnexpectedEOF instead of ErrLength.
var ErrLength = errors.New("bit: bit string length not a multiple of 8")

// InvalidByteError values describe errors resulting from an invalid byte in a bit string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("bit: invalid byte: %#U", rune(e))
}

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 8.
func EncodedLen(n int) int { return n * 8 }

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements bit encoding.
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
		dst[j+7] = bitTable[(v & b1)]
		j += 8
	}
	return len(src) * 8
}

// EncodeToString returns the bit encoding of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

// bufferSize is the number of bit characters to buffer in encoder and decoder.
const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns an io.Writer that writes bit characters to w.
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

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 8.
func DecodedLen(x int) int { return x / 8 }

// Decode decodes src into DecodedLen(len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// Decode expects that src contains only '0' or '1'
// characters and that src has multiple of 8 length.
// If the input is malformed, Decode returns the number
// of bytes decoded before the error.
func Decode(dst, src []byte) (int, error) {
	i, j := 0, 7
	for ; j < len(src); j += 8 {
		a, ok := fromBitChar(src[j-7])
		if !ok {
			return i, InvalidByteError(src[j-7])
		}
		b, ok := fromBitChar(src[j-6])
		if !ok {
			return i, InvalidByteError(src[j-6])
		}
		c, ok := fromBitChar(src[j-5])
		if !ok {
			return i, InvalidByteError(src[j-5])
		}
		d, ok := fromBitChar(src[j-4])
		if !ok {
			return i, InvalidByteError(src[j-4])
		}
		e, ok := fromBitChar(src[j-3])
		if !ok {
			return i, InvalidByteError(src[j-3])
		}
		f, ok := fromBitChar(src[j-2])
		if !ok {
			return i, InvalidByteError(src[j-2])
		}
		g, ok := fromBitChar(src[j-1])
		if !ok {
			return i, InvalidByteError(src[j-1])
		}
		h, ok := fromBitChar(src[j])
		if !ok {
			return i, InvalidByteError(src[j])
		}
		dst[i] = (a << 7) | (b << 6) | (c << 5) | (d << 4) | (e << 3) | (f << 2) | (g << 1) | h
		i++
	}
	if len(src)%8 != 0 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		for k := j - 7; k < (j-7)+len(src)%8; k++ {
			_, ok := fromBitChar(src[k])
			if !ok {
				return i, InvalidByteError(src[k])
			}
		}
		return i, ErrLength
	}
	return i, nil
}

// fromBitChar converts a bit character into its value and a success flag.
func fromBitChar(c byte) (byte, bool) {
	switch c {
	case '0':
		return 0, true
	case '1':
		return 1, true
	}

	return 0, false
}

// DecodeString returns the bytes represented by the bit string s.
//
// DecodeString expects that src contains only bit
// characters and that src has multiple of 8 length.
// If the input is malformed, DecodeString returns
// the bytes decoded before the error.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	// We can use the source slice itself as the destination
	// because the decode loop increments by one and then the 'seen' byte is not used anymore.
	n, err := Decode(src, src)
	return src[:n], err
}

type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

// NewDecoder returns an io.Reader that decodes bit characters from r.
// NewDecoder expects that r contain only an multiple of 8 length of bit characters.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < 8 && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in) // Copies 0 ~ 7 bytes
		numRead, d.err = d.r.Read(d.arr[numCopy:])
		d.in = d.arr[:numCopy+numRead]
		if d.err == io.EOF && len(d.in)%8 != 0 {
			if _, ok := fromBitChar(d.in[len(d.in)-1]); !ok {
				d.err = InvalidByteError(d.in[len(d.in)-1])
			} else {
				d.err = io.ErrUnexpectedEOF
			}
		}
	}

	// Decode internal buffer into output buffer
	if numAvail := len(d.in) / 8; len(p) > numAvail {
		p = p[:numAvail]
	}
	numDec, err := Decode(p, d.in[:len(p)*8])
	d.in = d.in[8*numDec:]
	if err != nil {
		d.in, d.err = nil, err // Decode error; discard input remainder
	}

	if len(d.in) < 8 {
		return numDec, d.err // Only expose errors when buffer fully consumed
	}
	return numDec, nil
}

// Dump returns a string that contains a bit dump of the given data. The format
// of the bit dump matches the output of `xxd -b` on the command line.
func Dump(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var buf strings.Builder
	// Dumper will write 72 bytes per complete 6 byte chunk, and at least
	// 67 bytes for whatever remains. Round the allocation up, since only a
	// maximum of 5 bytes will be wasted.
	buf.Grow((1 + ((len(data) - 1) / 6)) * 72)

	dumper := Dumper(&buf)
	dumper.Write(data)
	dumper.Close()
	return buf.String()
}

// Dumper returns a WriteCloser that writes a bit dump of all written data to
// w. The format of the dump matches the output of `xxd -b` on the command
// line.
func Dumper(w io.Writer) io.WriteCloser {
	return &dumper{w: w}
}

type dumper struct {
	w          io.Writer
	rightChars [7]byte
	buf        [14]byte
	used       int  // number of bytes in the current line
	n          uint // number of bytes, total
	closed     bool
}

func toChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}

func (h *dumper) Write(data []byte) (n int, err error) {
	if h.closed {
		return 0, errors.New("bit: dumper closed")
	}

	// Output lines look like:
	// 00000000: 00110000 00110001 00110010 00110011 00110100 00110101  012345
	// ^ offset                                                         ^ ASCII of line.
	for i := range data {
		if h.used == 0 {
			// At the beginning of a line we print the current
			// offset in bit.
			h.buf[0] = byte(h.n >> 24)
			h.buf[1] = byte(h.n >> 16)
			h.buf[2] = byte(h.n >> 8)
			h.buf[3] = byte(h.n)
			hex.Encode(h.buf[4:], h.buf[:4])
			h.buf[12] = ':'
			h.buf[13] = ' '
			_, err = h.w.Write(h.buf[4:])
			if err != nil {
				return
			}
		}
		Encode(h.buf[:], data[i:i+1])
		h.buf[8] = ' '
		l := 9
		if h.used == 5 {
			// There's an additional space after the 8th byte.
			h.buf[9] = ' '
			l = 10
		}
		_, err = h.w.Write(h.buf[:l])
		if err != nil {
			return
		}
		n++
		h.rightChars[h.used] = toChar(data[i])
		h.used++
		h.n++
		if h.used == 6 {
			h.rightChars[6] = '\n'
			_, err = h.w.Write(h.rightChars[:])
			if err != nil {
				return
			}
			h.used = 0
		}
	}
	return
}

func (h *dumper) Close() (err error) {
	// See the comments in Write() for the details of this format.
	if h.closed {
		return
	}
	h.closed = true
	if h.used == 0 {
		return
	}
	h.buf[0] = ' '
	h.buf[1] = ' '
	h.buf[2] = ' '
	h.buf[3] = ' '
	h.buf[4] = ' '
	h.buf[5] = ' '
	h.buf[6] = ' '
	h.buf[7] = ' '
	h.buf[8] = ' '
	h.buf[9] = ' '
	nBytes := h.used
	for h.used < 6 {
		l := 9
		if h.used == 5 {
			l = 10
		}
		_, err = h.w.Write(h.buf[:l])
		if err != nil {
			return
		}
		h.used++
	}
	h.rightChars[nBytes] = '\n'
	_, err = h.w.Write(h.rightChars[:nBytes+1])
	return
}
