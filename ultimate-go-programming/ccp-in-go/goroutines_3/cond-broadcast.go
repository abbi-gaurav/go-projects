package main

import (
	"fmt"
	"sync"
)

type button struct {
	clicked *sync.Cond
}

func main() {
	button := button{clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.clicked, func() {
		fmt.Println("maximize window")
		clickRegistered.Done()
	})
	subscribe(button.clicked, func() {
		fmt.Println("show diaglog box!")
		clickRegistered.Done()
	})
	subscribe(button.clicked, func() {
		fmt.Println("message flicked")
		clickRegistered.Done()
	})

	button.clicked.Broadcast()
	clickRegistered.Wait()
}
