package main

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	dfa := make(map[int]map[string][]string)
	a := make(map[string][]string)
	a["E"] = []string{".aA", ".bB"}
	a["A"] = []string{".a", "a.A"}
	dfa[0] = a
	b := make(map[string][]string)
	b["E"] = []string{".bA", "b.B"}
	b["A"] = []string{".b", "bB."}
	dfa[1] = b

	c := make(map[string][]string)
	c["E"] = []string{".bA", "b.B"}
	c["A"] = []string{".b", "bB."}
	dfa[2] = c
	idx := 2

	for i := 0; i < idx; i++ {
		isSame := true
		if len(dfa[i]) != len(dfa[idx]) {
			continue
		}
		for sKey, list := range dfa[idx] {
			if _, ok := dfa[idx][sKey]; ok {
				if len(list) != len(dfa[idx][sKey]) {
					isSame = false
					break
				}
				for j := range dfa[idx][sKey] {
					if dfa[idx][sKey][j] != dfa[i][sKey][j] {
						isSame = false
						break
					}
				}
			} else {
				isSame = false
			}
			if !isSame {
				break
			}
		}
		if isSame {
			fmt.Println("YES")
			return
		}
	}

}
