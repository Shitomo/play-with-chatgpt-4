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
		log.Print("handle hello")
		resp, err := http.Get(os.Getenv("API_SERVER_URL") + "/hello")
		if err != nil {
			log.Print(err)
			fmt.Fprintf(w, "error %v+", err)

			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "error %v+", err)

			return
		}
		fmt.Fprintf(w, string(body)+",Goodbye World")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
