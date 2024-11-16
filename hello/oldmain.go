package hello

import (
	"bufio"
	"fmt"
	"os"
	// "strings"
)

// "hello"

func main() {
	// fmt.Println(hello.Say(os.Args[1:]))
	// hello.ReadFromStdInAndCalcAverage()
	// arrStuff()

	// otherArrStuff()
	byteStuff()
	// 10001001 first byte
	// 11011001 second byte

	//     1010 1010
	//       0xAA
	// AND 0011 0010
	//       0x32
	//     0010 0010

	//     1010 1010
	// OR  0011 0010
	//     1011 1010

}

func byteStuff() {
	var firstByte uint8 = 0xFF
	fmt.Printf("%08b \n", firstByte)

	var mask uint8 = 252 // 1111 1100 aka 0xFC
	fmt.Printf("%08b \n", mask)

	var otherMask uint8 = 0x4 // 0000 0100

	//     1010 0111
	// AND 0000 0100
	//     0000 0100
	var only1bit uint8 = firstByte & otherMask
	fmt.Printf("%08b \n", only1bit)

	aMap := map[byte]string{
		0: "AL",
		1: "CL",
		2: "DL",
	}
}

func otherArrStuff() {
	t := []byte("string")
	fmt.Println(len(t), t)
	fmt.Println(t[len(t)-2:])
}

func arrStuff() {
	var c = []int{1, 2, 3}
	c = append(c, 4, 5)
	c = c[2:4]
	for i := 0; i < len(c); i++ {
		fmt.Println(c[i])
	}

}

func stringStuff2() {
	args := os.Args

	if len(args) < 3 {
		i, err := fmt.Fprintln(os.Stderr, "Failed, bro, not enough args")
		fmt.Println(i)
		fmt.Println(err)
		os.Exit(-1)
	}

	// old, new := os.Args[1], os.Args[2]
	old := os.Args[1]
	new := os.Args[2]
	fmt.Printf("%s %s", old, new)

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		// s := strings.Split(scan.Text(), old)
		// t := strings.Join(s, new)
		fmt.Printf(scan.Text())
	}
}

func stringStuff() {
	fr := "élitè"
	fmt.Printf("%8T %[1]v %d\n", fr, len(fr))
	fmt.Printf("%8T %[1]v %d\n", []rune(fr), len([]rune(fr)))
	fmt.Printf("%8T %[1]v %d\n", []byte(fr), len([]byte(fr)))
}
