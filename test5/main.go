package main

/*
Task:
Given a string, return the first non-repeating character.

input: "swiss"
output: "w"
*/
func main() {
	println(nonRepeating("swiss"))   // Output: "w"
	println(nonRepeating("racecar")) // Output: "e"
}

func nonRepeating(s string) string {
	charCount := make(map[rune]int)

	// Count occurrences of each character
	for _, char := range s {
		charCount[char]++
	}

	// Find the first non-repeating character
	for _, char := range s {
		if charCount[char] == 1 {
			return string(char)
		}
	}

	return ""
}
