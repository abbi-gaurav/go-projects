package main

import (
	"fmt"
)

type data struct {
	name string
	age  int
}

func (d data) displayName() {
	fmt.Println("Name is ", d.name)
}

func (d *data) setAge(age int) {
	d.age = age
	fmt.Println(d.name," Age is: ", d.age)
}

func main() {
	d := data{
		name: "ss",
	}
	d.displayName()
	d.setAge(33)

	//what go does underneath
	data.displayName(d)
	(* data).setAge(&d, 44)

	f1 := d.displayName
	f1()

	d.name = "ga"
	f1()

	f2 := d.setAge
	d.name = "dd"
	f2(45)
}
