package main

import "fmt"

func main() {
	var strings [2]string
	strings[0] = "apple"
	strings[1] = "banana"

	for i, fruit := range strings {
		fmt.Println(i, fruit)
	}
}
