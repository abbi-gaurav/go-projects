package main

import "fmt"

type notifier interface {
	notify()
}

type userd1 struct {
	name  string
	email string
}

func (u *userd1) notify() {
	fmt.Printf("Sending userd1 email to %s<%s>\n", u.name, u.email)
}

type admin struct {
	userd1
	level string
}

func (a *admin) notify() {
	fmt.Printf("Sending admin email to %s<%s>\n", a.name, a.email)
}

func sendNotification(n notifier) {
	n.notify()
}

func main() {
	ad := admin{
		userd1: userd1{
			name:  "ss",
			email: "s@e",
		},
		level: "super",
	}

	sendNotification(&ad)

	ad.userd1.notify()

	ad.notify()
}
