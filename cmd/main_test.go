package main

import "testing"

func TestDecodeInstruction_SingleRegister(t *testing.T) {
	var byte1 byte = 0b10001001
	var byte2 byte = 0b11011001

	output := DecodeInstructions([]byte{byte1, byte2})
	expected := "mov cx, bx\n"

	if output != expected {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func TestDecodeInstruction_MultipleRegisters(t *testing.T) {
	var byte1 byte = 0b10001001
	var byte2 byte = 0b11011001

	output := DecodeInstructions([]byte{byte1, byte2, byte1, byte2})
	expected := "mov cx, bx\nmov cx, bx\n"

	if output != expected {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}
