package string_processing

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"log"
)

func lineListSource(ctx context.Context, lines ...string) (<-chan string, <-chan error, error) {
	if len(lines) == 0 {
		return nil, nil, errors.New("no line provided")
	}

	out := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)

		for index, line := range lines {
			if line == "" {
				errc <- errors.New(fmt.Sprintf("line %v is empty", index+1))
				return
			}
			select {
			case out <- line:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

func xtf(ctx context.Context, base int, in <-chan string) (<-chan int64, <-chan error, error) {
	if base < 2 {
		return nil, nil, errors.New(fmt.Sprintf("Invalid base %v", base))
	}

	out := make(chan int64)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for line := range in {
			n, err := strconv.ParseInt(line, base, 64)
			if err != nil {
				errc <- err
				return
			}
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc, nil
}

func sink(ctx context.Context, in <-chan int64) (<-chan error, error) {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		for n := range in {
			if n >= 100 {
				errc <- errors.New(fmt.Sprintf("number %v is too large", n))
				return
			}
			fmt.Printf("sink : %v\n", n)
		}
	}()
	return errc, nil
}

func RunSimplePipeline(base int, lines []string) error {
	ctx, canceFunc := context.WithCancel(context.Background())
	defer canceFunc()

	var errcList []<-chan error

	linec, errc, err := lineListSource(ctx, lines...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	numberc, errc, err := xtf(ctx, base, linec)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	errc, err = sink(ctx, numberc)
	if err != nil {
		return err
	}

	fmt.Println("pipeline started. Waiting for pipeline to complete.")
	err = waitForPipeline(errcList...)

	if err != nil{
		log.Panic(err)
	}

	return err
}

func waitForPipeline(errorChannels ...<-chan error) error {
	errc := mergeErrorChannels(errorChannels...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil

}

func mergeErrorChannels(errorChannels ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(errorChannels))

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(errorChannels))

	for _, c := range errorChannels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func splitter(ctx context.Context, in <-chan int64) (<-chan int64, <-chan int64, <-chan error, error) {
	out1 := make(chan int64)
	out2 := make(chan int64)
	errc := make(chan error, 1)

	go func() {
		defer close(out1)
		defer close(out2)
		defer close(errc)

		for n := range in {
			select {
			case out1 <- n:
			case <-ctx.Done():
				return
			}
			select {
			case out2 <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out1, out2, errc, nil
}

func squarer(ctx context.Context, in <-chan int64) (<-chan int64, <-chan error, error) {
	out := make(chan int64)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc, nil
}

