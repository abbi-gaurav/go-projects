package main

import "fmt"

type user struct {
	name string
}

type notifiable interface {
	notify()
}

func (u user) notify() {
	fmt.Println("userData ", u.name, "notified")
}

type Pipe struct {
	name string
}

func (p *Pipe) notify() {
	fmt.Println("Pipe ", p.name, " notified")
}

func doNotify(notifiable notifiable) {
	notifiable.notify()
}

func main() {
	u := user{"ss"}
	p := Pipe{"gg"}
	p.notify()

	doNotify(u)
	//does not compile
	//doNotify(p)
}
