package main

import (
	"enigmacamp.com/golatihanlagi/delivery"
	_ "github.com/lib/pq"
)

func main() {
	delivery.NewServer().Run()
}
