package main

import (
	"net/http"
)

func main() {

	serv := New()

	http.ListenAndServe(":8080", serv.Routes())

}
