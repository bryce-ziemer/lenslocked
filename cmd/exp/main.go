package main

import (
	"fmt"
)

func blah() error {
	panic("...")
}

func main() {
	numbers := []int{1, 2, 3}
	fmt.Println(numbers[4])
}
