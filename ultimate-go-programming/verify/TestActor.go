package main

import "github.com/go-learning-projects/ultimate-go-programming/concurrency/actors"
import "time"

var actor *actors.MyActor

func main() {
	actor = actors.CreateActor()
	go incrementBy1()
	go incrementBy2()

	time.Sleep(10 * time.Second)

}

func incrementBy2() {
	actor.Increment(2)
}

func incrementBy1() {
	actor.Increment(1)
}
