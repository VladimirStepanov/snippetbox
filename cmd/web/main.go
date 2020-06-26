package main

import "flag"

func main() {
	addr := flag.String("addr", ":8080", "Listen addr")

	flag.Parse()
	serv := New(*addr)
	serv.Start()

}
