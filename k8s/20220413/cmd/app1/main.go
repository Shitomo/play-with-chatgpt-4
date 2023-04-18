package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handle hello")
		fmt.Fprintf(w, "Hello World")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
