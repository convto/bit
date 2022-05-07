# bit
This is a bit encoding library, supporting encoding and decoding.  
It also provides io.Reader and io.Writer.  
The implementation is based on encoding/hex of the standard package.

## What is this?
This package provides Encode/Decode functionality for byte sequences into bit strings.

For example, a byte represented by the hexadecimal number `ff` would be encoded as `1111111111` .

## Why is this necessary?
Go does not (as far as I can tell) have the flexibility to output raw byte sequences as bit strings.  
This can be a problem in log output when, for example, parsing a binary message fails.

Bit output with padding like `fmt.Sprintf("%08b", buf)` is close,  
but I created this package for more flexible handling (e.g. io.Reader and io.Writer support).

## Usage

Only basic Encode/Decode is queried here. Please refer to the example for details.

### Encode

```go
src := []byte("Hello Gopher!")

dst := make([]byte, EncodedLen(len(src)))
Encode(dst, src)

fmt.Printf("%s\n", dst)

// Output:
// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
}
```

### Decode

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

## LICENSE
MIT