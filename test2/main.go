package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
	goal:
	- decode strings with the following pattern: k[encoded_string], where the encoded_string inside the square brackets is repeated exactly k times.
	- nested patterns are allowed, e.g., 2[ab3[c5d]2x2[i]] should be decoded to abcd5d5d5dxxiiabcd5d5d5dxxii
*/

func main() {
	str := "2[ab3[c5d]2x2[i]]"
	println("RESULT: ", decodeString(str))
}

func decodeString(s string) string {
	previous := s
	for {
		s = closing(0, s)
		if previous == s {
			break
		}
		previous = s
	}
	return s
}

func closing(i int, s string) string {
	for k := i + 1; k < len(s); k++ {
		if s[k] == ']' {
			count, err := strconv.Atoi(string(s[i-1]))
			if err != nil {
				panic(err)
			}
			str := s[i+1 : k]
			var res string
			for range count {
				res += str
			}

			replaceString := fmt.Sprintf("%d[%s]", count, str)
			fmt.Printf("replacing %s TO %s \n", replaceString, res)
			return strings.Replace(s, replaceString, res, -1)
		}
		if s[k] == '[' {
			return closing(k, s)
		}
	}
	return s
}
