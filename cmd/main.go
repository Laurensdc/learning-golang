package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

var debugging bool = false

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

	for bytePointer := 0; bytePointer < len(bytes); {
		if debugging {
			fmt.Printf("Processing byte %v\n", bytePointer)
		}

		var byte1 byte = bytes[bytePointer]
		var byte2 byte = bytes[bytePointer+1]

		if byte1&0b1111_1100>>2 == 0b100010 {
			// Register/memory to/from register

			d := byte1 & 0b0000_0010 >> 1 // 0: source is in reg field, 1: dest is in reg field
			w := byte1 & 0b0000_0001

			mod := byte2 & 0b1100_0000 >> 6
			const (
				MemoryModeNoDisplacement    = iota // mod == 0b00
				MemoryMode8BitDisplacement         // mod == 0b01
				MemoryMode16BitDisplacement        // mod == 0b10
				RegisterMode                       // mod == 0b11
			)

			reg := byte2 & 0b0011_1000 >> 3 // name of register
			rm := byte2 & 0b0000_0111       // also name of register, or maybe name of memory Register/Memory R/M

			if mod == MemoryModeNoDisplacement {
				if rm == 0b110 {
					// if r/m is 110, we have 16 bit displacement (exception case lol)
					// FIXME: Unhandled instruction
					bytePointer += 4
				} else {
					decodedInstruction += memoryModeNoDisplacement(d, w, reg, rm)
					bytePointer += 2
				}
			} else if mod == MemoryMode8BitDisplacement {
				byte3 := bytes[bytePointer+2]
				decodedInstruction += memoryMode8BitDisplacement(d, w, reg, rm, byte3)
				bytePointer += 3
			} else if mod == MemoryMode16BitDisplacement {
				byte3 := bytes[bytePointer+2]
				byte4 := bytes[bytePointer+3]
				decodedInstruction += memoryMode16BitDisplacement(d, w, reg, rm, byte3, byte4)
				bytePointer += 4
			} else if mod == RegisterMode {
				decodedInstruction += registerMode(d, w, reg, rm)
				bytePointer += 2
			}
		} else if byte1&0b1111_0000>>4 == 0b1011 {
			// immediate to register
			w := byte1 & 0b0000_1000 >> 3
			reg := byte1 & 0b0000_0111

			if w == 0 {
				// not wide, read 1 byte (signed) as immediate value
				dataValueStr := fmt.Sprintf("%v", int8(byte2))
				decodedInstruction += "mov " + decodeRegister(w, reg) + ", " + dataValueStr + "\n"
				bytePointer += 2
			} else if w == 1 {
				// "w" for wide, read 2 bytes as immediate value
				byte3 := bytes[bytePointer+2]
				dataValueStr := bytesToStr([2]byte{byte2, byte3})
				decodedInstruction += "mov " + decodeRegister(w, reg) + ", " + dataValueStr + "\n"
				bytePointer += 3
			}
		} else {
			fmt.Println("Unspecified operation")
			os.Exit(1)
		}
	}

	if debugging {
		fmt.Printf("Decoded output omg\n%v\n", decodedInstruction)
	}

	return decodedInstruction
}

// mov instruction in register mode
func registerMode(d, w, reg, rm byte) string {
	regStr := decodeRegister(w, reg)
	rmStr := decodeRegister(w, rm)
	operands := orderOperands(d, regStr, rmStr)

	return "mov " + operands + "\n"
}

// mov instruction in memory mode with no displacement
func memoryModeNoDisplacement(d, w, reg, rm byte) string {
	regStr := decodeRegister(w, reg)
	effectiveAddress := decodeEffectiveAddress(rm, "")
	operands := orderOperands(d, regStr, effectiveAddress)

	return "mov " + operands + "\n"
}

// mov instruction in memory mode with 8 bit displacement
func memoryMode8BitDisplacement(d, w, reg, rm, byte3 byte) string {
	displacement := fmt.Sprintf("%v", byte3)
	regStr := decodeRegister(w, reg)
	effectiveAddress := decodeEffectiveAddress(rm, displacement)
	operands := orderOperands(d, regStr, effectiveAddress)

	return "mov " + operands + "\n"
}

// mov instruction in memory mode with 16 bit displacement
func memoryMode16BitDisplacement(d, w, reg, rm, byte3, byte4 byte) string {
	displacement := bytesToStr([2]byte{byte3, byte4})
	regStr := decodeRegister(w, reg)
	effectiveAddress := decodeEffectiveAddress(rm, displacement)
	operands := orderOperands(d, regStr, effectiveAddress)

	return "mov " + operands + "\n"
}

// Put the bytes in the right order and interpret as unsigned int
func bytesToStr(bytes [2]byte) string {
	dataValue := int16(binary.LittleEndian.Uint16(bytes[:]))

	dataValueStr := fmt.Sprintf("%v", dataValue)
	return dataValueStr
}

func decodeRegister(w, reg byte) string {
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

	return mapRegister[w][reg]
}

func decodeEffectiveAddress(rm byte, displacement string) string {
	mapEffectiveAddress := map[byte]string{
		0b000: "bx + si",
		0b001: "bx + di",
		0b010: "bp + si",
		0b011: "bp + di",
		0b100: "si",
		0b101: "di",
		0b110: "bp",
		0b111: "bx",
	}

	var displacementStr string
	if displacement == "" || displacement == "0" {
		displacementStr = ""
	} else {
		displacementStr = " + " + displacement
	}

	return "[" + mapEffectiveAddress[rm] + displacementStr + "]"
}

// Order both operands based on the value of d
func orderOperands(d byte, reg, regOrMemoryAddress string) string {
	if d == 0 {
		// 0: source is in reg field
		return regOrMemoryAddress + ", " + reg
	} else if d == 1 {
		// 1: destination is in reg field
		return reg + ", " + regOrMemoryAddress
	} else {
		fmt.Printf("Received invalid value %v for d\n", d)
		os.Exit(1)
		return ""
	}
}
