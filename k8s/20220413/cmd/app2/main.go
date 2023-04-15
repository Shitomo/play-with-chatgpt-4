package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(os.Getenv("ENDPOINT") + "/hello")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(body)+"Goodbye World")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
