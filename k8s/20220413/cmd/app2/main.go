package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	app1v1 "github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1"
	"github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1/app1v1connect"
	"github.com/bufbuild/connect-go"
)

func main() {

	url := os.Getenv("API_SERVER_URL")
	if url == "" {
		log.Fatal("API_SERVER_URL is not set")
	}

	client := app1v1connect.NewHelloServiceClient(http.DefaultClient, url)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Print("handle hello")

		res, err := client.Hello(r.Context(), connect.NewRequest(&app1v1.HelloRequest{}))
		if err != nil {
			log.Printf("failed to call hello: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s, and Goodbye", res.Msg.Message)))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
