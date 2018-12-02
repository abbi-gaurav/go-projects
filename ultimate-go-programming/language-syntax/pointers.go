package main

import "fmt"

type user struct {
	name  string
	email string
}

func main() {
	lecture2()
}

func lecture2() {
	fmt.Printf("%+v", stayOnStack())
	println("----")
	fmt.Printf("%+v",*escapeToHeap())
}
func escapeToHeap() *user {
	u := user{
		name:  "ga",
		email: "ga@ga.com",
	}
	return &u
}
func stayOnStack() user {
	u := user{
		name:  "ga",
		email: "ga@ga.com",
	}
	return u
}

func lecture1() {
	count := 10
	println("before", "count - ", count, "address - ", &count)
	increment(count)
	println("after", "count - ", count, "address - ", &count)
	increment2(&count)
	println("after2", "count - ", count, "address - ", &count)
}
func lecture4()  {

}

func increment(number int) {
	number++
	println("number", number, "address", &number)
}

func increment2(numberPointer *int) {
	*numberPointer++
	println("numberPointer", *numberPointer, &numberPointer, numberPointer)
}
