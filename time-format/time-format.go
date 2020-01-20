package main

import (
	"log"
	"time"
)

func main() {
	parse, err := time.Parse(time.RFC3339, "2019-05-06T01:13:19.533")
	log.Println(err)
	log.Println(parse)
}
