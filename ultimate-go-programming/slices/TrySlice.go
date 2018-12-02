package main

import "fmt"

func main() {
	sharing()
}

func s1() {
	slice := make([]string, 2)
	slice[0] = "apple"
	slice[1] = "banana"
	//slice[3] = "RE"

	fmt.Println(slice)
}

func sharing() {
	slice := make([] int, 3)
	slice[0] = 1
	slice[1] = 1
	slice[2] = 1

	sharedPtr := &slice[1]
	slice = append(slice, 2)
	slice[1]++

	fmt.Println("Pointer:", *sharedPtr, "element", slice[1])
}

func sliceUnchanged() {
	slice1 := make([]string, 4, 6)
	slice1[0] = "apple"
	slice1[1] = "banana"
	slice1[2] = "grape"
	slice1[3] = "orange"

	slice2 := slice1[2:3:3]
	fmt.Println(slice1)
	fmt.Println(slice2)

	slice2 = append(slice2, "new")
	fmt.Println(slice1)
	fmt.Println(slice2)
}

func sliceChanged() {
	slice1 := make([]string, 4)
	slice1[0] = "apple"
	slice1[1] = "banana"
	slice1[2] = "grape"
	slice1[3] = "orange"

	slice2 := slice1[2:4]
	fmt.Println(slice1)
	fmt.Println(slice2)

	slice2[1] = "changed"
	fmt.Println(slice1)
	fmt.Println(slice2)
}

func s2() {
	var data []string

	lastCap := cap(data)

	for record := 1; record <= 102400; record++ {
		data = append(data, fmt.Sprintf("Rec: %d", record))

		if lastCap != cap(data) {
			capChange := float64(cap(data)-lastCap) / float64(lastCap) * 100

			lastCap = cap(data)

			fmt.Printf("\nIndex[%d] cap[%d - %f]\n", len(data), lastCap, capChange)
		}
	}
}
