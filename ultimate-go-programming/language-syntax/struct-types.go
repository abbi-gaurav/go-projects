package main

import "fmt"

type example struct {
	flag    bool
	counter int16
	pi      float32
}

type a struct {
	same bool
}

type b struct {
	same bool
}

func main() {
	var e1 example //this is 8 bytes;

	fmt.Printf("%+v", e1)

	e2 := example{
		flag:    true,
		counter: 10,
		pi:      3.141592,
	}
	fmt.Println()
	fmt.Println("Flag", e2.flag)
	fmt.Println("Counter", e2.counter)
	fmt.Println("pi", e2.pi)

	e3 := struct {
		length int32
	}{
		length: 38,
	}

	fmt.Printf("%+v\n", e3)

	var a1 a
	var b1 b

	a1 = a(b1)
	fmt.Println(a1, b1)

	unnamed := struct {
		same bool
	}{
		same: true,
	}

	a1 = unnamed
	fmt.Println(a1, unnamed)
}
