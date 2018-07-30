package main

import (
	"fmt"
	"time"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	fmt.Println(time.Now().UnixNano())
	if err := fasthttp.ListenAndServe(":8080", fastHTTPHandler); err != nil {
		log.Fatalf("Error in listenAndServe: %s", err)
	}
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi!")
}