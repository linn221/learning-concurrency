package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Hour * 1)
		w.Write([]byte("one hour has passed"))
	})
	http.HandleFunc("/cancel-me", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		go func() {
			<-ctx.Done()
			fmt.Println("request has been cancelled")
		}()
		for i := range 10 {
			time.Sleep(time.Second * 1)
			fmt.Printf("%d second has passed\n", i+1)
		}
		_, err := w.Write([]byte("10 second has been passed"))
		fmt.Println(err.Error())
	})

	http.ListenAndServe(":8080", nil) // uses the default http server
}
