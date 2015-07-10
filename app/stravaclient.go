package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "kombagger ", 0)

const baseAPI = "https://www.strava.com/api/v3/"
const segmentExploreAPI = "segments/explore"
const segmentLeaderboardAPI = "segments/%v/leaderboard"
const retrieveActivityAPI = "activities/%v"

const accessToken = "a22c2bb5f6270300e5fd562c702084126301e286"
const authParam = "access_token=" + accessToken

type HttpClient interface {
	Get(url string) []byte
}

type DefaultHttpClient struct{}

func (defaultHttpClient DefaultHttpClient) Get(url string) []byte {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return bytes
}

var httpClient HttpClient = DefaultHttpClient{}

// allow test double
func setHttpClient(httpClient_ HttpClient) {
	httpClient = httpClient_
}

type Point struct {
	lat  float32
	long float32
}

type BoundingBox struct {
	topRight   Point
	bottomLeft Point
}

type Activity struct {
	Id             int             `json:"id"`
	SegmentEfforts []SegmentEffort `json:"segment_efforts"`
}

type SegmentEffort struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Segment Segment `json:"segment"`
}

// deserialization helper
type segments struct {
	Segments []Segment
}

type Segment struct {
	Id           int
	Name         string
	Distance     float32
	End_Latlng   []float32
	Start_Latlng []float32
}

func (segment *Segment) Start() Point {
	return Point{segment.Start_Latlng[0], segment.Start_Latlng[1]}
}

func (segment *Segment) End() Point {
	return Point{segment.End_Latlng[0], segment.End_Latlng[1]}
}

type Effort struct {
	Id             int    `json:"effort_id"`
	AthleteName    string `json:"athlete_name"`
	AthleteId      int    `json:"athlete_id"`
	ElapsedTime    int    `json:"elapsed_time"`
	MovingTime     int    `json:"moving_time"`
	Rank           int    `json:"rank"`
	AthleteProfile string `json:"athlete_profile"`
}

// deserialization helper
type efforts struct {
	Efforts []Effort `json:"entries"`
}

func parseErrorString(resp []byte) string {
	return fmt.Sprintf("unable to parse json response \n %s", resp)
}

// GetSegments returns a selection of segments within the boundingBox specified.
func GetSegments(boundingBox BoundingBox) []Segment {

	boundsParam := "bounds=" + fmt.Sprintf("%v", boundingBox.bottomLeft.lat) + "," + fmt.Sprintf("%v", boundingBox.bottomLeft.long) + "," + fmt.Sprintf("%v", boundingBox.topRight.lat) + "," + fmt.Sprintf("%v", boundingBox.topRight.long)

	url := baseAPI + segmentExploreAPI + "?" + boundsParam + "&" + authParam

	logger.Println("calling Strava API: " + url)
	resp := httpClient.Get(url)

	segments := new(segments)
	err := json.Unmarshal(resp, segments)

	if err != nil {
		panic(parseErrorString(resp))
	}

	return segments.Segments

}

func GetSegmentLeaderboard(segmentId int) []Effort {

	url := fmt.Sprintf(baseAPI+segmentLeaderboardAPI+"?"+authParam, segmentId)

	resp := httpClient.Get(url)

	logger.Println(string(resp))

	efforts := new(efforts)
	err := json.Unmarshal(resp, efforts)

	if err != nil {
		panic(parseErrorString(resp))
	}

	return efforts.Efforts
}

func GetActivity(activityId int) (Activity, error) {

	url := fmt.Sprintf(baseAPI+retrieveActivityAPI+"?"+authParam, activityId)

	logger.Println("calling Strava API: " + url)
	resp := httpClient.Get(url)

	activity := Activity{}
	err := json.Unmarshal(resp, &activity)

	if err != nil {
		logger.Printf("%v", err)
		return Activity{}, errors.New("error calling Strava API")
	} else {
		return activity, nil
	}

}
