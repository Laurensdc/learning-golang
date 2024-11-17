package main

import (
	"fmt"
	"os"
)

var debugging bool = true

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		debugging = true
	}

	bytes, err := os.ReadFile("../binaries/listing_0038_many_register_mov")

	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	if debugging {
		fmt.Printf("Read file %08b", bytes)
	}

	instruction := DecodeInstructions(bytes)

	fmt.Println(instruction)
}

func DecodeInstructions(bytes []byte) string {

	decodedInstruction := ""

	for i := 0; i < len(bytes); {
		var byte1 byte = bytes[i]
		var byte2 byte = bytes[i+1]

		if debugging {
			fmt.Printf("Reading bytes %b %b\n", byte1, byte2)
		}

		if byte1&0b1111_1100>>2 == 0b100010 {
			//register/memory to/from register

			d := byte1 & 0b0000_0010 >> 1 // 0: source is in reg field, 1: dest is in reg field
			w := byte1 & 0b0000_0001

			mod := byte2 & 0b1100_0000 >> 6
			reg := byte2 & 0b0011_1000 >> 3 // name of register
			rm := byte2 & 0b0000_0111       // also name of register, or maybe name of memory Register/Memory R/M

			if mod == 0b00 {
				// Everything is encoded in these 2 bytes

				if rm == 0b110 {
					// if r/m is 110, we have 16 bit displacement (exception case lol)
					i += 4
				} else {
					// no displacement
					i += 2
				}
			}
			if mod == 0b11 {
				// Everything is encoded in these 2 bytes

				// when we do			mov ax, bx
				// it means				ax = bx
				orderedOperands := map[byte]string{
					0: decodeOperands(w, rm, reg),
					1: decodeOperands(w, reg, rm),
				}

				decodedInstruction += "mov " + orderedOperands[d] + "\n"

				i += 2
			}
			if mod == 0b01 {
				// Read 3 bytes
				i += 3
			}
			if mod == 0b10 {
				// Read 4 bytes
				i += 4
			}
		} else if byte1&0b1111_0000>>4 == 0b1011 {
			// immediate to register
			w := byte1 & 0b0000_1000 >> 3

			if w == 0 {
				i += 2
			} else if w == 1 {
				i += 3
			}
		} else {
		}
	}

	if debugging {
		fmt.Printf("Decoded output omg\n%v\n", decodedInstruction)
	}

	return decodedInstruction
}

func decodeOperands(w, reg1, reg2 byte) string {
	mapRegister := map[byte]map[byte]string{
		0: {
			0b000: "al",
			0b001: "cl",
			0b010: "dl",
			0b011: "bl",
			0b100: "ah",
			0b101: "ch",
			0b110: "dh",
			0b111: "bh",
		},
		1: {
			0b000: "ax",
			0b001: "cx",
			0b010: "dx",
			0b011: "bx",
			0b100: "sp",
			0b101: "bp",
			0b110: "si",
			0b111: "di",
		},
	}

	return mapRegister[w][reg1] + ", " + mapRegister[w][reg2]
}
