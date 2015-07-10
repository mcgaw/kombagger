package app

import (
	"errors"
	"time"
)

var referenceActivities = []int{

	// round the loch anti-clockwise
	341750344,
	// rond the lock clockwise
	// TODO
	// brackletter spean reverse
	339664883,
	// regular spean loop
	326016101,
	// dam ride with brakcletter then glen nevis
	310099490,
	// lochailort
	300447396,
	// liddesdale regular
	310099490,
	// liddesdale reverse
	// TODO

}

var fortWilliamBoundingBox = BoundingBox{Point{56.8295301, -5.0756684}, Point{56.80433499999999, -5.1304722}}

func init() {

	go wait()
}

// hacky, wait for the DB to initialize itself
func wait() {
	time.Sleep(time.Second * 10)
	go pollStrava()
}

func pollStrava() {

	logger.Println("updating Strava data")

	segments, err := compileSegmentsList()

	if err == nil {
		for _, segment := range segments {
			leaderboard := GetSegmentLeaderboard(segment.Id)

			logger.Println("adding new KOM to the database")
			AddKOM(leaderboard[0])
		}
		logger.Println("finished updating the Strava data")
	} else {
		logger.Printf("there was an error updating the Strava data: %v", err)
	}

	time.Sleep(time.Minute * 60)

}

func compileSegmentsList() ([]Segment, error) {

	segmentMap := make(map[int]Segment)

	for _, activityId := range referenceActivities {
		activity, err := GetActivity(activityId)
		if err != nil {
			return []Segment{}, errors.New("error getting activity data from Strava API")
		}
		for _, segmentEffort := range activity.SegmentEfforts {
			segmentMap[segmentEffort.Segment.Id] = segmentEffort.Segment
		}
	}

	var segments []Segment

	for _, segment := range segmentMap {
		segments = append(segments, segment)
	}

	return segments, nil

}
