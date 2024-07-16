package gogit

import (
	"fmt"
	"testing"
)

func AssertEq(t *testing.T, actual any, expected any) {
	if expected != actual {
		fmt.Printf("Failed! expected: [%v] actual: [%v]\n", expected, actual)
		t.Fail()
	}
}
