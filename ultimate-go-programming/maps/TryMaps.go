package main

import "fmt"

type user struct {
	firstName string
	lastName  string
}

func main() {
	s1()
}

func s1() {
	users := make(map[string]user)
	users["g"] = user{"g", "a"}
	users["s"] = user{"s", "s"}

	for key, value := range users {
		fmt.Println(key, value)
	}

	delete(users, "g")

	u, found := users["g"]
	fmt.Println(u,found)
}

func s2() {
	users := map[string]user{
		"g": user{"g", "a"},
		"s": user{"s", "s"},
	}

	for key, value := range users {
		fmt.Println(key, value)
	}
}
