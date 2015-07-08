package app

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type LeaderboardController struct{}

func (controller *LeaderboardController) GetLeaderboard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	leaderboard := GetRiderLeaderboard(1)
	json, err := json.MarshalIndent(leaderboard, "", "    ")

	if err != nil {
		w.Header().Set("Status", "500")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	fmt.Fprint(w, string(json))

}

func NewLeaderboardController() *LeaderboardController {
	return &LeaderboardController{}
}
