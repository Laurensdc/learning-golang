package main

import (
	"fmt"
	"os"
)

var debugging bool

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		debugging = true
	}

	bytes, err := os.ReadFile("../binaries/listing_0038_many_register_mov")

	if debugging {
		fmt.Printf("Read file %v", bytes)
	}

	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	// Every thing is 16 bits, so I need to fetch them
	// 16 bits at a time
	if len(bytes)%2 != 0 {
		fmt.Printf("Didn't provide 16 bit instructions, cannot decode\n")
		os.Exit(1)
	}

	instruction := DecodeInstructions(bytes)

	fmt.Println(instruction)
}

func DecodeInstructions(bytes []byte) string {
	decodedInstruction := ""

	// Each set of 2 bytes is 1 instruction
	for i := 0; i < len(bytes)-1; i += 2 {
		var byte1 byte = bytes[i]
		var byte2 byte = bytes[i+1]

		if debugging {
			fmt.Printf("Reading bytes %b %b\n", byte1, byte2)
		}

		var opcode = byte1 & 0b1111_1100 >> 2 // first 6 bytes is opcode
		var d = byte1 & 0b0000_0010 >> 1      // 0: source is in reg field, 1: dest is in reg field
		var w = byte1 & 0b0000_0001

		var mod = byte2 & 0b1100_0000 >> 6
		var reg = byte2 & 0b0011_1000 >> 3 // name of register
		var rm = byte2 & 0b0000_0111       // also name of register, or maybe name of memory Register/Memory R/M

		if debugging {
			fmt.Printf("all the stuff: opcode %b d %b w %b mod %b reg %b rm %b\n", opcode, d, w, mod, reg, rm)
		}

		if opcode == 0b100010 { // mov
			decodedInstruction += "mov "
			// d is 0 so
			// source is in reg field

			// when we do mov ax, bx
			// it means   mov dest, source
			// like						ax = bx

			// so now we know that source is in reg field
			// so bx is in reg field
			// that means that the dest is in the rm field

			// now I need to get dest first

			if d == 0 { // source is in reg field, dest is in r/m field
				if w == 1 {
					// Map for when W == 1
					registerMapW1 := map[byte]string{
						0b000: "ax",
						0b001: "cx",
						0b010: "dx",
						0b011: "bx",
						0b100: "sp",
						0b101: "bp",
						0b110: "si",
						0b111: "di",
					}

					destRegister := registerMapW1[rm]
					decodedInstruction += destRegister + ", "

					sourceRegister := registerMapW1[reg]
					decodedInstruction += sourceRegister
				} else if w == 0 {
					// TODO: Create map for w == 0
				}

			}

		}

		decodedInstruction += "\n"
	}

	if debugging {
		fmt.Printf("Decoded output omg\n%v\n", decodedInstruction)

	}
	return decodedInstruction
}
