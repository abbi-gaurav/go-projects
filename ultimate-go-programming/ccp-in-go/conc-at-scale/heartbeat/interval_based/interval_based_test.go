package interval_based

import (
	"testing"
	"time"
)

func doWork(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	hbs := make(chan interface{}, 1)
	ints := make(chan int)

	go func() {
		defer close(hbs)
		defer close(ints)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case hbs <- struct{}{}:
					default:
					}
				case ints <- n:
					continue numLoop
				}
			}
		}
	}()
	return hbs, ints
}

func TestDoWork_GenerateAllNums(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{1, 2, 3, 4, 5, 6, 7}
	const timeout = 2 * time.Second

	hbs, results := doWork(done, timeout/2, intSlice...)

	<-hbs

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v, expected %v, but received %v", i, expected, r)
			}
			i++
		case <-hbs:
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
