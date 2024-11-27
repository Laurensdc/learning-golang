package hello

import (
	"fmt"
	"os"
	"strings"
)

func Say(names []string) string {
	if len(names) == 0 {
		return ("Hello, Laurens")
	}

	return ("Hello, " + strings.Join(names, ", "))
}

func ReadFromStdInAndCalcAverage() {
	var sum float64
	var n int

	for {
		var val float64

		_, err := fmt.Fscanln(os.Stdin, &val)

		if err != nil {
			break
		}

		sum += val
		n++
	}

	if n == 0 {
		os.Exit(-1)

	}

	fmt.Println("The average is ", sum/float64(n))
}
