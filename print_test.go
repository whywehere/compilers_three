package main

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	var s string
	fmt.Scanln("hello world", &s)
	println(s)
}
