package main

import (
	"fmt"
	"go-class/decoder"
	"os"
)

var debugging bool = true

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		debugging = true
	}

	bytes, err := os.ReadFile("binaries/listing_0041_add_sub_cmp_jnz")

	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	if debugging {
		fmt.Printf("Read file %#08b", bytes)
	}

	instruction := decoder.DecodeInstructions(bytes)

	fmt.Println(instruction)
}
