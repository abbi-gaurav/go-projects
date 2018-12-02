package main

import "fmt"

type userData struct {
	name  string
	email string
}

//fmt.Stringer interface
func (u *userData) String() string {
	return fmt.Sprintf("My name is %q and my email is %q", u.name, u.email)
}

func main() {
	ud := userData{
		name:  "SS",
		email: "se",
	}

	fmt.Println(ud)

	//overrides
	fmt.Println(&ud)
}
