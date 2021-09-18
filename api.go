package main

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func GetUpcomingRacesHandler(ctx *fasthttp.RequestCtx) {
	race, err := GetNextRace()
	if err != nil {
		if err == ErrNoRace {
			NoRaceHandler(ctx)
			return
		}
		ctx.Error("", 500)
	}

	json, err := json.Marshal(race)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(json)
	return
}

func GetCurrentRaceHandler(ctx *fasthttp.RequestCtx) {
	race, err := GetCurrentRace()
	if err != nil {
		if err == ErrNoRace {
			NoRaceHandler(ctx)
			return
		}
		ctx.Error("", 500)
	}

	json, err := json.Marshal(race)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(json)
	return
}

func NoRaceHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(`{"status": "NO RACE"}`)
	return
}
