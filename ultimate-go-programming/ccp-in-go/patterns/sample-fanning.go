package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/ccp-in-go/patterns/fanning"
	"runtime"
	"time"
)

type item struct {
	id            int
	name          string
	packingEffort time.Duration
}

func prepareItems(done <-chan interface{}) <-chan interface{} {
	items := make(chan interface{})
	itemsToShip := []item{
		{0, "Shirt", 1 * time.Second},
		{1, "Legos", 1 * time.Second},
		{2, "TV", 5 * time.Second},
		{3, "Bananas", 1 * time.Second},
		{4, "Hat", 1 * time.Second},
		{5, "Phone", 2 * time.Second},
		{6, "Plates", 3 * time.Second},
		{7, "Computer", 5 * time.Second},
		{8, "Pint Glass", 3 * time.Second},
		{9, "Watch", 2 * time.Second},
	}

	go func() {
		defer close(items)
		for _, item := range itemsToShip {
			select {
			case <-done:
				return
			case items <- item:
			}
		}
	}()

	return items
}

func packItems(done <-chan interface{}, items <-chan interface{}, workerId int) <-chan interface{} {
	packages := make(chan interface{})
	go func() {
		defer close(packages)

		for i := range items {
			itemObj := i.(item)
			select {
			case <-done:
				return
			case packages <- itemObj.id:
				time.Sleep(itemObj.packingEffort)
				fmt.Printf("worker #%d: shipping package no. %d, tood %ds to pack\n", workerId, itemObj.id, itemObj.packingEffort/time.Second)
			}
		}
	}()
	return packages
}

func main() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	items := prepareItems(done)

	workers := fanning.FanOut(done, packItems, runtime.NumCPU(), items)

	numPackages := 0

	fannedIn := fanning.FanIn(done, workers...)

	for range fannedIn {
		numPackages++
	}

	fmt.Printf("took %fs to ship %d packages\n", time.Since(start).Seconds(), numPackages)

}
