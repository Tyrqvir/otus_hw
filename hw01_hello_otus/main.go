package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	sourceString := "Hello, OTUS!"

	reverseString := stringutil.Reverse(sourceString)
	fmt.Print(reverseString)
}
