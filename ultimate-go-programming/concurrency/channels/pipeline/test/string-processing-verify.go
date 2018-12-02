package main

import "github.com/go-learning-projects/ultimate-go-programming/concurrency/channels/pipeline/string-processing"

func main() {
	var strings []string
	strings = append(strings, "1")

	string_processing.RunSimplePipeline(2, strings)
}
