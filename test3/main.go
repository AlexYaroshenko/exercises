package main

import (
	"fmt"
)

/*
	goal:
	- assign seats to passengers strictly following the order of preference:

	  1. The middle seat (if the number of seats is odd, the exact middle; if even, the left middle)
	  2. The seat to the left of the middle seat
	  3. The seat to the right of the middle seat
	  4. The seat two to the left of the middle seat
	  5. The seat two to the right of the middle seat
	  ... and so on, alternating left and right until all seats are filled.
	- once all seats are filled, the process repeats for the next group of passengers.
	- return a slice of integers representing the order of seat assignments for the given number of passengers and seats.
*/

func places(passengers, sits int) []int {
	var order []int
	occupied := make(map[int]struct{}, sits)
	people := 0
	for i := 0; i < passengers; i++ {
		sit := findSit(occupied, sits)
		order = append(order, sit)
		occupied[sit] = struct{}{}
		people++
		if people == sits {
			people = 0
			occupied = make(map[int]struct{}, sits)
		}
	}

	return order
}

func findSit(occupied map[int]struct{}, sits int) int {
	var sit int
	if sits%2 == 0 {
		sit = sits / 2
	} else {
		sit = sits/2 + 1
	}

	left, right := sit, sit
	for {
		if _, ok := occupied[sit]; !ok {
			return sit
		}
		left--
		if _, ok := occupied[left]; !ok {
			if left != 0 {
				return left
			}
		}
		right++
		if _, ok := occupied[right]; !ok {
			return right
		}
	}
}

func main() {
	fmt.Println(places(6, 5))
	fmt.Println(places(9, 6))
}
