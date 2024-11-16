package main

import "testing"

func TestDecodeInstruction(t *testing.T) {
	var byte1 byte = 0b10001001
	var byte2 byte = 0b11011001

	output := DecodeInstruction([2]byte{byte1, byte2})
	expected := "mov cx, bx"

	if output != expected {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}
