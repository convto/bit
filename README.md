# bit
[![Go Reference](https://pkg.go.dev/badge/github.com/convto/bit.svg)](https://pkg.go.dev/github.com/convto/bit) [![Go Report Card](https://goreportcard.com/badge/github.com/convto/bit)](https://goreportcard.com/report/github.com/convto/bit) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This package provides encoding for converting byte sequences to bit strings.  
For example, a byte represented by the hexadecimal number `ff` would be encoded as `11111111` .

Inspired by the standard package encoding/hex. I also got some implementation hints from it.

## Why is this necessary?
Byte sequences can be output as bit strings using the standard package `fmt` as follows

```go
for i := 0; i < len(src); i++ {
    fmt.Printf("%08b", src[i])
}
```

In some cases, this is sufficient. However, this does not implement `io.Reader` or `io.Writer` , so flexible handling such as stream support is not possible.  
In addition, it is tedious to output in a human-readable format like `xxd -b`.

Therefore, outputting bit strings using the standard package `fmt` is inconvenient, for example, when debugging a program that evaluates binaries.

I created this package for more flexible handling (e.g. `io.Reader` and `io.Writer` support, Or `Dump()` output support like `xxd -b` ).

## Usage

Here are the basics. If you want more details, please refer to [example_test](./bit_example_test.go) or [package documentation](https://pkg.go.dev/github.com/convto/bit).

### Encode

Encode the given byte sequences.

```go
src := []byte("Hello Gopher!")

dst := make([]byte, EncodedLen(len(src)))
Encode(dst, src)

fmt.Printf("%s\n", dst)

// Output:
// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
```

### Decode

Decode takes as input a bit-encoded byte sequences.  
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

Dump returns output like `xxd -b` .

```go
dump := Dump([]byte("dump test"))
fmt.Printf("%s\n", dump)

// Output:
// 00000000: 01100100 01110101 01101101 01110000 00100000 01110100  dump t
// 00000006: 01100101 01110011 01110100                             est
```

## LICENSE
MIT
