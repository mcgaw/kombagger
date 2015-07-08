package app

import (
	"time"
)

var fortWilliamBoundingBox = BoundingBox{Point{56.8295301, -5.0756684}, Point{56.80433499999999, -5.1304722}}

func init() {

	go wait()
}

func wait() {
	time.Sleep(time.Second * 10)
	go pollStrava()
}

func pollStrava() {

	logger.Println("Calling the Strava API")

	segments := GetSegments(fortWilliamBoundingBox)
	for _, segment := range segments {
		leaderboard := GetSegmentLeaderboard(segment.Id)

		logger.Println("adding new KOM to the database")
		AddKOM(leaderboard[0])
	}

	time.Sleep(time.Minute * 1000)

}
