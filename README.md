# bit
[![Go Reference](https://pkg.go.dev/badge/github.com/convto/bit.svg)](https://pkg.go.dev/github.com/convto/bit) [![Go Report Card](https://goreportcard.com/badge/github.com/convto/bit)](https://goreportcard.com/report/github.com/convto/bit) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This package provides Encode/Decode functionality for byte sequences into bit strings.  
For example, a byte represented by the hexadecimal number `ff` would be encoded as `1111111111` .

Inspired by the standard package encoding/hex. I also got some implementation hints from it.

## Why is this necessary?
Go does not (as far as I can tell) have the flexibility to output raw byte sequences as bit strings.  
This can be a problem in log output when, for example, parsing a binary message fails.

Bit output with padding like `fmt.Sprintf("%08b", buf)` is close,  
but I created this package for more flexible handling (e.g. `io.Reader` and `io.Writer` support, Or `Dump()` output support like `xxd -b` ).

## Usage

Here are the basics. If you want more details, please refer to example_test or the package documentation.

### Encode

Encode the given byte sequence.

```go
src := []byte("Hello Gopher!")

dst := make([]byte, EncodedLen(len(src)))
Encode(dst, src)

fmt.Printf("%s\n", dst)

// Output:
// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
```

### Decode

The input to the decoding process is a bit-encoded byte sequence.  
Eight characters represent one byte, so the input must be a multiple of 8 bytes.

```go
src := []byte("01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001")

dst := make([]byte, DecodedLen(len(src)))
n, err := Decode(dst, src)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("%s\n", dst[:n])

// Output:
// Hello Gopher!
```

### Dump

Dump returns the same output as `xxd -b` .

```go
dump := Dump([]byte("dump test"))
fmt.Printf("%s\n", dump)

// Output:
// 00000000: 01100100 01110101 01101101 01110000 00100000 01110100  dump t
// 00000006: 01100101 01110011 01110100                             est
```

## LICENSE
MIT