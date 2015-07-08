package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mcgaw/kombagger/app"
	"net/http"
)

func main() {
	r := httprouter.New()
	c := app.NewLeaderboardController()

	r.GET("/leaderboard", c.GetLeaderboard)

	http.ListenAndServe("localhost:9090", r)
}
