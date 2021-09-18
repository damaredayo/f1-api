package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	// Update calander at start so first API request is fast
	updateCalander()

	h := fasthttp.CompressHandler(handle)
	if err := fasthttp.ListenAndServe(":8080", h); err != nil {
		log.Fatalf(err.Error())
	}

}
