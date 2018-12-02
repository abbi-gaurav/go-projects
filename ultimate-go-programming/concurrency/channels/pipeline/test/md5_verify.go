package main

import (
	"os"
	"fmt"
	"sort"
	"github.com/go-learning-projects/ultimate-go-programming/concurrency/channels/pipeline/md5/parallel"
	"github.com/go-learning-projects/ultimate-go-programming/concurrency/channels/pipeline/md5/serial"
	"github.com/go-learning-projects/ultimate-go-programming/concurrency/channels/pipeline/md5/bounded"
)

func main() {
	computeType := os.Args[2]
	switch computeType {
	case "serial":
		serial.MD5All(os.Args[1])
	case "parallel":
		parallel.MD5All(os.Args[1])
	case "bounded":
		bounded.MD5All(os.Args[1])
	}
	m, err := parallel.MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var paths []string

	for path := range m {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
	fmt.Println("compute -- ", computeType)
}
