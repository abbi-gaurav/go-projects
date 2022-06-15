package main

import "fmt"

type Subtractable interface {
	~int | int32 | float64 | float32
}

type Results[T Subtractable] []T

type MyOwnInteger int

type Movable[S Subtractable] interface {
	move(S)
}

type Person[S Subtractable] struct {
	name string
}

func (p Person[S]) move(meters S) {
	fmt.Printf("%s moved %d meters\n", p.name, meters)
}

func main() {
	fmt.Println(subtract(1, 0))
	fmt.Println(subtract(12.4, 13.8))
	fmt.Println(subtract2(14.4, 13.8))
	fmt.Println(subtract2[int](2, 3))

	var myInt MyOwnInteger
	myInt = 10
	fmt.Println(subtract2(myInt, 9))

	var resultStorage Results[int]
	resultStorage = append(resultStorage, subtract[int](5, 8))
	fmt.Println("ResultStorage:", resultStorage)

	p := Person[int]{name: "John"}
	fmt.Println(move(p, 10, 2))
}

func move[S Subtractable, V Movable[S]](v V, distance S, meters S) S {
	v.move(meters)
	return subtract2(distance, meters)
}

func subtract[V int64 | int | float32 | float64](a, b V) V {
	return a - b
}

func subtract2[V Subtractable](a, b V) V {
	return a - b
}
