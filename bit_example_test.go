package bit

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func ExampleEncode() {
	src := []byte("Hello Gopher!")

	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)

	fmt.Printf("%s\n", dst)

	// Output:
	// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
}

func ExampleEncodeToString() {
	src := []byte("Hello Gopher!")
	encodedStr := EncodeToString(src)

	fmt.Printf("%s\n", encodedStr)

	// Output:
	// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
}

func ExampleEncoder_Write() {
	src := []byte("Hello Gopher!")
	enc := NewEncoder(os.Stdout)
	enc.Write(src)
	// Output:
	// 01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001
}

func ExampleDecode() {
	src := []byte("01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001")

	dst := make([]byte, DecodedLen(len(src)))
	n, err := Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", dst[:n])

	// Output:
	// Hello Gopher!
}

func ExampleDecodeString() {
	const s = "01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001"
	decoded, err := DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", decoded)

	// Output:
	// Hello Gopher!
}

func ExampleDecoder_Read() {
	src := []byte("01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001")
	buf := bytes.NewBuffer(src)
	dec := NewDecoder(buf)
	dst := make([]byte, DecodedLen(len(src)))
	dec.Read(dst)
	fmt.Printf("%s\n", dst)

	// Output:
	// Hello Gopher!
}
