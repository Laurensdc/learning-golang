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
	bytes := []byte{137, 217, 136, 229, 137, 218, 137, 222, 137, 251, 136, 200, 136, 237, 137, 195, 137, 243, 137, 252, 137, 197}
	output := DecodeInstructions(bytes)
	expected :=
		`mov cx, bx
mov ch, ah
mov dx, bx
mov si, bx
mov bx, di
mov al, cl
mov ch, ch
mov bx, ax
mov bx, si
mov sp, di
mov bp, ax
`

	if output != expected {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func TestDecodeManyMore(t *testing.T) {
	input := []byte{0b10001001, 0b11011110, 0b10001000, 0b11000110, 0b10110001, 0b00001100, 0b10110101, 0b11110100, 0b10111001, 0b00001100, 0b00000000, 0b10111001, 0b11110100, 0b11111111, 0b10111010, 0b01101100, 0b00001111, 0b10111010, 0b10010100, 0b11110000, 0b10001010, 0b00000000, 0b10001011, 0b00011011, 0b10001011, 0b01010110, 0b00000000, 0b10001010, 0b01100000, 0b00000100, 0b10001010, 0b10000000, 0b10000111, 0b00010011, 0b10001001, 0b00001001, 0b10001000, 0b00001010, 0b10001000, 0b01101110, 0b00000000}
	expected := DecodeInstructions(input)
	output := `mov si, bx
mov dh, al
mov cl, 12
mov ch, -12
mov cx, 12
mov cx, -12
mov dx, 3948
mov dx, -3948
mov al, [bx + si]
mov bx, [bp + di]
mov dx, [bp]
mov ah, [bx + si + 4]
mov al, [bx + si + 4999]
mov [bx + di], cx
mov [bp + si], cl
mov [bp], ch`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}
