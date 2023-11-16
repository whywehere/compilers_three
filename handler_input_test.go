package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestHandlerInput(t *testing.T) {
	handlerGrammar := make(map[string][]string)

	for _, grammar := range input {

		// 取出E->aA|bB 中的 aA|bB
		str := strings.Split(grammar, "->") // fmt.Println(rightStr)
		leftStr := str[0]

		rightStr := str[1]
		// 以'|' 分割aA|bB
		rightGrammars := strings.Split(rightStr, "|")
		for _, rg := range rightGrammars {
			handlerGrammar[leftStr] = append(handlerGrammar[leftStr], rg)
		}
	}
	// debug print
	for k, v := range handlerGrammar {
		fmt.Println(k, " ", v)
	}
}
