package main

import (
	"flag"
	"fmt"
)

func main() {
	addr := flag.String("addr", ":8080", "Listen addr")

	flag.Parse()
	serv := New(*addr)

	if err := serv.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
