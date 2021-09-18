package main

import "github.com/valyala/fasthttp"

func handle(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("", fasthttp.StatusMethodNotAllowed)
		return
	}

	switch string(ctx.Path()) {
	case "/api/getupcoming":
		GetUpcomingRacesHandler(ctx)
	case "/api/getcurrent":
		GetCurrentRaceHandler(ctx)
	case "/":
		// Just a simple status check
		ctx.Success("text/plain", []byte{})
	}
}
