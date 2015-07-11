package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mcgaw/kombagger/app"
	"net/http"
	"os"
)

func main() {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false

	c := app.NewLeaderboardController()
	r.GET("/api/leaderboard", c.GetLeaderboard)

	r.NotFound = http.FileServer(http.Dir("../ember-kombagger/dist"))

	http.ListenAndServe("localhost:"+os.Getenv("PORT"), r)
}
