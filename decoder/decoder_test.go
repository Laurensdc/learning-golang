package decoder

import (
	"testing"
)

func Test_Single_Reg_To_Reg(t *testing.T) {
	var byte1 byte = 0b10001001
	var byte2 byte = 0b11011001

	output := DecodeInstructions([]byte{byte1, byte2})
	expected := "mov cx, bx\n"

	if output != expected {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func Test_RegMemToFromReg_RegisterMode_11(t *testing.T) {
	bytes := []byte{0b10001001, 0b11011001, 0b10001000, 0b11100101, 0b10001001, 0b11011010, 0b10001001, 0b11011110, 0b10001001, 0b11111011, 0b10001000, 0b11001000, 0b10001000, 0b11101101, 0b10001001, 0b11000011, 0b10001001, 0b11110011, 0b10001001, 0b11111100, 0b10001001, 0b11000101}
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

func Test_RegMemToFromReg_MemoryMode_00(t *testing.T) {
	input := []byte{0b10001010, 0b00000000, 0b10001011, 0b00011011}
	output := DecodeInstructions(input)
	expected := `mov al, [bx + si]
mov bx, [bp + di]
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func Test_RegMemToFromReg_MemoryMode_01_8BitDisplacement(t *testing.T) {
	input := []byte{0b10001011, 0b01010110, 0b00000000, 0b10001010, 0b01100000, 0b00000100}
	output := DecodeInstructions(input)
	expected :=
		`mov dx, [bp]
mov ah, [bx + si + 4]
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func Test_RegMemToFromReg_MemoryMode_10_16BitDisplacement(t *testing.T) {
	input := []byte{0b10001010, 0b10000000, 0b10000111, 0b00010011}
	output := DecodeInstructions(input)
	expected := `mov al, [bx + si + 4999]
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}

}

func Test_RegMemToFromReg_MemoryMode_10_16BitDisplacement_DFlag(t *testing.T) {
	input := []byte{0b10001001, 0b00001001, 0b10001000, 0b00001010, 0b10001000, 0b01101110, 0b00000000}
	output := DecodeInstructions(input)
	expected := `mov [bx + di], cx
mov [bp + si], cl
mov [bp], ch
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func Test_ImmediateToRegister(t *testing.T) {
	input := []byte{0b10110001, 0b00001100, 0b10110101, 0b11110100, 0b10111001, 0b00001100, 0b00000000, 0b10111001, 0b11110100, 0b11111111, 0b10111010, 0b01101100, 0b00001111, 0b10111010, 0b10010100, 0b11110000}
	output := DecodeInstructions(input)
	expected := `mov cl, 12
mov ch, -12
mov cx, 12
mov cx, -12
mov dx, 3948
mov dx, -3948
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func TestDecodeMovs(t *testing.T) {
	input := []byte{0b10001001, 0b11011110, 0b10001000, 0b11000110, 0b10110001, 0b00001100, 0b10110101, 0b11110100, 0b10111001, 0b00001100, 0b00000000, 0b10111001, 0b11110100, 0b11111111, 0b10111010, 0b01101100, 0b00001111, 0b10111010, 0b10010100, 0b11110000, 0b10001010, 0b00000000, 0b10001011, 0b00011011, 0b10001011, 0b01010110, 0b00000000, 0b10001010, 0b01100000, 0b00000100, 0b10001010, 0b10000000, 0b10000111, 0b00010011, 0b10001001, 0b00001001, 0b10001000, 0b00001010, 0b10001000, 0b01101110, 0b00000000}
	output := DecodeInstructions(input)
	expected := `mov si, bx
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
mov [bp], ch
`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}

func TestAddSubMovCmpJump(t *testing.T) {
	t.Skip("Skipping the big bang test because yeah it's too much at once")

	input := []byte{
		// reg/mem with reg to either
		0b00000011, 0b00011000, 0b00000011, 0b01011110, 0b00000000,
		// immediate to reg/mem
		// add si, 2...
		0b10000011, 0b11000110, 0b00000010,
		0b10000011, 0b11000101, 0b00000010,
		0b10000011, 0b11000001, 0b00001000,

		// reg/mem with register to either + displacement
		0b00000011, 0b01011110, 0b00000000, 0b00000011, 0b01001111, 0b00000010,
		0b00000010, 0b01111010, 0b00000100, 0b00000011, 0b01111011, 0b00000110, 0b00000001, 0b00011000, 0b00000001, 0b01011110,
		0b00000000, 0b00000001, 0b01011110, 0b00000000, 0b00000001, 0b01001111, 0b00000010, 0b00000000, 0b01111010, 0b00000100,
		0b00000001, 0b01111011, 0b00000110, 0b10000000, 0b00000111, 0b00100010, 0b10000011, 0b10000010, 0b11101000, 0b00000011,
		0b00011101, 0b00000011, 0b01000110, 0b00000000, 0b00000010, 0b00000000, 0b00000001, 0b11011000, 0b00000000, 0b11100000,
		0b00000101, 0b11101000, 0b00000011, 0b00000100, 0b11100010, 0b00000100, 0b00001001, 0b00101011, 0b00011000, 0b00101011,
		0b01011110, 0b00000000, 0b10000011, 0b11101110, 0b00000010, 0b10000011, 0b11101101, 0b00000010, 0b10000011, 0b11101001,
		0b00001000, 0b00101011, 0b01011110, 0b00000000, 0b00101011, 0b01001111, 0b00000010, 0b00101010, 0b01111010, 0b00000100,
		0b00101011, 0b01111011, 0b00000110, 0b00101001, 0b00011000, 0b00101001, 0b01011110, 0b00000000, 0b00101001, 0b01011110,
		0b00000000, 0b00101001, 0b01001111, 0b00000010, 0b00101000, 0b01111010, 0b00000100, 0b00101001, 0b01111011, 0b00000110,
		0b10000000, 0b00101111, 0b00100010, 0b10000011, 0b00101001, 0b00011101, 0b00101011, 0b01000110, 0b00000000, 0b00101010,
		0b00000000, 0b00101001, 0b11011000, 0b00101000, 0b11100000, 0b00101101, 0b11101000, 0b00000011, 0b00101100, 0b11100010,
		0b00101100, 0b00001001, 0b00111011, 0b00011000, 0b00111011, 0b01011110, 0b00000000, 0b10000011, 0b11111110, 0b00000010,
		0b10000011, 0b11111101, 0b00000010, 0b10000011, 0b11111001, 0b00001000, 0b00111011, 0b01011110, 0b00000000, 0b00111011,
		0b01001111, 0b00000010, 0b00111010, 0b01111010, 0b00000100, 0b00111011, 0b01111011, 0b00000110, 0b00111001, 0b00011000,
		0b00111001, 0b01011110, 0b00000000, 0b00111001, 0b01011110, 0b00000000, 0b00111001, 0b01001111, 0b00000010, 0b00111000,
		0b01111010, 0b00000100, 0b00111001, 0b01111011, 0b00000110, 0b10000000, 0b00111111, 0b00100010, 0b10000011, 0b00111110,
		0b11100010, 0b00010010, 0b00011101, 0b00111011, 0b01000110, 0b00000000, 0b00111010, 0b00000000, 0b00111001, 0b11011000,
		0b00111000, 0b11100000, 0b00111101, 0b11101000, 0b00000011, 0b00111100, 0b11100010, 0b00111100, 0b00001001, 0b01110101,
		0b00000010, 0b01110101, 0b11111100, 0b01110101, 0b11111010, 0b01110101, 0b11111100, 0b01110100, 0b11111110, 0b01111100,
		0b11111100, 0b01111110, 0b11111010, 0b01110010, 0b11111000, 0b01110110, 0b11110110, 0b01111010, 0b11110100, 0b01110000,
		0b11110010, 0b01111000, 0b11110000, 0b01110101, 0b11101110, 0b01111101, 0b11101100, 0b01111111, 0b11101010, 0b01110011,
		0b11101000, 0b01110111, 0b11100110, 0b01111011, 0b11100100, 0b01110001, 0b11100010, 0b01111001, 0b11100000, 0b11100010,
		0b11011110, 0b11100001, 0b11011100, 0b11100000, 0b11011010, 0b11100011, 0b11011000,
	}

	output := DecodeInstructions(input)
	expected := `add bx, [bx+si]
	add bx, [bp]
	add si, 2
	add bp, 2
	add cx, 8
	add bx, [bp + 0]
	add cx, [bx + 2]
	add bh, [bp + si + 4]
	add di, [bp + di + 6]
	add [bx+si], bx
	add [bp], bx
	add [bp + 0], bx
	add [bx + 2], cx
	add [bp + si + 4], bh
	add [bp + di + 6], di
	add byte [bx], 34
	add word [bp + si + 1000], 29
	add ax, [bp]
	add al, [bx + si]
	add ax, bx
	add al, ah
	add ax, 1000
	add al, -30
	add al, 9
	sub bx, [bx+si]
	sub bx, [bp]
	sub si, 2
	sub bp, 2
	sub cx, 8
	sub bx, [bp + 0]
	sub cx, [bx + 2]
	sub bh, [bp + si + 4]
	sub di, [bp + di + 6]
	sub [bx+si], bx
	sub [bp], bx
	sub [bp + 0], bx
	sub [bx + 2], cx
	sub [bp + si + 4], bh
	sub [bp + di + 6], di
	sub byte [bx], 34
	sub word [bx + di], 29
	sub ax, [bp]
	sub al, [bx + si]
	sub ax, bx
	sub al, ah
	sub ax, 1000
	sub al, -30
	sub al, 9
	cmp bx, [bx+si]
	cmp bx, [bp]
	cmp si, 2
	cmp bp, 2
	cmp cx, 8
	cmp bx, [bp + 0]
	cmp cx, [bx + 2]
	cmp bh, [bp + si + 4]
	cmp di, [bp + di + 6]
	cmp [bx+si], bx
	cmp [bp], bx
	cmp [bp + 0], bx
	cmp [bx + 2], cx
	cmp [bp + si + 4], bh
	cmp [bp + di + 6], di
	cmp byte [bx], 34
	cmp word [4834], 29
	cmp ax, [bp]
	cmp al, [bx + si]
	cmp ax, bx
	cmp al, ah
	cmp ax, 1000
	cmp al, -30
	cmp al, 9
	test_label0:
	jnz test_label1
	jnz test_label0
	test_label1:
	jnz test_label0
	jnz test_label1
	label:
	je label
	jl label
	jle label
	jb label
	jbe label
	jp label
	jo label
	js label
	jne label
	jnl label
	jg label
	jnb label
	ja label
	jnp label
	jno label
	jns label
	loop label
	loopz label
	loopnz label
	jcxz label
	`

	if expected != output {
		t.Errorf("Expected\n%v\ngot\n%v\n", expected, output)
	}
}
