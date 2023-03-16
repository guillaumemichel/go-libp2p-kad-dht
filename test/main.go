package main

// test program. we can use the go one.
import (
	"encoding/binary" // varint is here
	"fmt"
)

func main() {
	ints := []uint64{1, 127, 128, 255, 300, 500, 501, 16384}
	for _, i := range ints {
		buf := make([]byte, 10)
		n := binary.PutUvarint(buf, uint64(i))

		hexStr := fmt.Sprintf("%x", i)
		if len(hexStr)%2 == 1 {
			hexStr = "0" + hexStr
		}
		fmt.Print(i, " (0x"+hexStr+")\t=> ")
		for c := 0; c < n; c++ {
			fmt.Printf("%08b ", int(buf[c]))
		}
		fmt.Printf("(0x%x)\n", buf[:n])
	}
}
