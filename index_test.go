package main

import (
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	s := ".BB"
	index := strings.Index(s, ".")
	str := s[:index] + string(s[index+1]) + "." + s[index+2:]
	println(str)
}
