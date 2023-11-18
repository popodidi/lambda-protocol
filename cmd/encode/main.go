package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	input := `{"a":1,"b":2}`
	hexStr := hex.EncodeToString([]byte(input))
	fmt.Println(hexStr)
}
