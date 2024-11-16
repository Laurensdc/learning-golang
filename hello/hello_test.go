package hello

import (
	"testing"
)

func TestSayHelloOneName(t *testing.T) {
	want := "Hello, Josh"
	got := Say([]string{"Josh"})

	if want != got {
		t.Errorf("wanted %s -- got %s", want, got)
	}
}

func TestSayHelloNoNames(t *testing.T) {
	want := "Hello, Laurens"
	got := Say([]string{})

	if want != got {
		t.Errorf("wanted %s -- got %s", want, got)
	}
}

func TestSayHelloMultipleNames(t *testing.T) {
	want := "Hello, Josh, Paula, Ron"
	got := Say([]string{"Josh", "Paula", "Ron"})

	if want != got {
		t.Errorf("wanted %s -- got %s", want, got)
	}
}
