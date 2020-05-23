package main

import "net/http"

func main() {
	err := http.ListenAndServe(":3000", http.FileServer(http.Dir("./web")))
	if err != nil {
		panic(err)
	}
}
