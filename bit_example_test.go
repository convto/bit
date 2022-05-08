package bit

import (
	"bytes"
	"fmt"
	"io"
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

func ExampleNewEncoder() {
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

func ExampleNewDecoder() {
	src := []byte("01001000011001010110110001101100011011110010000001000111011011110111000001101000011001010111001000100001")

	buf := bytes.NewBuffer(src)
	dec := NewDecoder(buf)

	io.Copy(os.Stdout, dec)

	// Output:
	// Hello Gopher!
}

func ExampleDump() {
	dump := Dump([]byte("dump test"))
	fmt.Printf("%s\n", dump)

	// Output:
	// 00000000: 01100100 01110101 01101101 01110000 00100000 01110100  dump t
	// 00000006: 01100101 01110011 01110100                             est
}

func ExampleDumper() {
	d := Dumper(os.Stdout)
	d.Write([]byte("dump test"))
	d.Close()

	// Output:
	// 00000000: 01100100 01110101 01101101 01110000 00100000 01110100  dump t
	// 00000006: 01100101 01110011 01110100                             est
}
