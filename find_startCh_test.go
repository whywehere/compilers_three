package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestFindStartCh(t *testing.T) {
	for _, r := range leftList {
		isStart := true
		for _, grammar := range handlerGrammar {
			for _, s := range grammar {
				if strings.Contains(s, r) {
					isStart = false
					break
				}
			}
		}
		if isStart {
			startChar = r
		}
	}
	fmt.Println(startChar)
}
