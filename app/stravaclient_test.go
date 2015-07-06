package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockHttpClient struct {
	url          string
	responseBody []byte
}

func (client *MockHttpClient) Get(url string) []byte {
	client.url = url
	return client.responseBody
}

var testSegment = `
{
  "segments": [
    {
      "id": 229781,
      "name": "Hawk Hill",
      "climb_category": 1,
      "climb_category_desc": "4",
      "avg_grade": 5.7,
      "start_latlng": [
        37.8331119,
        -122.4834356
      ],
      "end_latlng": [
        37.8280722,
        -122.4981393
      ],
      "elev_difference": 152.8,
      "distance": 2684.8,
      "points": "}g|eFnm@n@Op@VJr@"
    }
  ]
}`

var testSegments = `
{
  "segments": [
    {
      "id": 229781,
      "name": "Hawk Hill",
      "climb_category": 1,
      "climb_category_desc": "4",
      "avg_grade": 5.7,
      "start_latlng": [
        37.8331119,
        -122.4834356
      ],
      "end_latlng": [
        37.8280722,
        -122.4981393
      ],
      "elev_difference": 152.8,
      "distance": 2684.8,
      "points": "}g|eFnm@n@Op@VJr@"
    },
    {
      "id": 632535,
      "name": "Hawk Hill Upper Conzelman to Summit",
      "climb_category": 0,
      "climb_category_desc": "NC",
      "avg_grade": 8.10913,
      "start_latlng": [
        37.8334451,
        -122.4941994
      ],
      "end_latlng": [
        37.8281297,
        -122.4980005
      ],
      "elev_difference": 67.29200000000003,
      "distance": 829.834,
      "points": "_j|eFvc@p@SbAu@h@Qn@?RTDH"
    }
  ]
}
`

var testSegmentLeaderboard = `
  {
  "entry_count": 7037,
  "entries": [
    {
      "athlete_name": "Jim Whimpey",
      "athlete_id": 123529,
      "athlete_gender": "M",
      "average_hr": 190.5,
      "average_watts": 460.8,
      "distance": 2659.89,
      "elapsed_time": 370,
      "moving_time": 360,
      "start_date": "2013-03-29T13:49:35Z",
      "start_date_local": "2013-03-29T06:49:35Z",
      "activity_id": 46320211,
      "effort_id": 801006623,
      "rank": 1,
      "athlete_profile": "http://pics.com/227615/large.jpg"
    },
    {
      "athlete_name": "Chris Zappala",
      "athlete_id": 11673,
      "athlete_gender": "M",
      "average_hr": null,
      "average_watts": 368.3,
      "distance": 2705.7,
      "elapsed_time": 374,
      "moving_time": 374,
      "start_date": "2012-02-23T14:50:16Z",
      "start_date_local": "2012-02-23T06:50:16Z",
      "activity_id": 4431903,
      "effort_id": 83383918,
      "rank": 2,
      "athlete_profile": "http://pics.com/227615/large.jpg"
    }
  ]
}`

func Test_GetSegments(t *testing.T) {

	mockHttpClient := MockHttpClient{responseBody: []byte(testSegment)}

	setHttpClient(&mockHttpClient)

	testBox := BoundingBox{
		Point{50.0, 0.0},
		Point{0.0, 50.0}}

	segments := GetSegments(testBox)

	assert.Equal(t, "https://www.strava.com/api/v3/"+"segments/explore"+"?bounds=0,50,50,0&access_token=a22c2bb5f6270300e5fd562c702084126301e286",
		mockHttpClient.url, "incorrect request url")

	assert.Equal(t, 1, len(segments), "wrong number of segments returned")

	assert.Equal(t, 229781, segments[0].Id, "incorrect id")
	assert.Equal(t, "Hawk Hill", segments[0].Name, "incorrect name")
	assert.Equal(t, "2684.8", fmt.Sprintf("%v", segments[0].Distance), "incorrect distance")
	assert.Equal(t, Point{37.8331119, -122.4834356}, segments[0].Start(), "incorrect start point")
	assert.Equal(t, Point{37.8280722, -122.4981393}, segments[0].End(), "incorrect end point")

}

func Test_GetSegmentLeaderboard(t *testing.T) {

	mockHttpClient := MockHttpClient{responseBody: []byte(testSegmentLeaderboard)}

	setHttpClient(&mockHttpClient)

	leaderboard := GetSegmentLeaderboard(1)

	assert.Equal(t, 2, len(leaderboard), "incorret number of efforts")

	// should be ordered
	for i, effort := range leaderboard {
		assert.Equal(t, i+1, effort.Rank)
	}

	effort := leaderboard[0]

	assert.Equal(t, "Jim Whimpey", effort.AthleteName, "incorrect athelete name")
	assert.Equal(t, 123529, effort.AthleteId, "incorrect athelete id")
	assert.Equal(t, 801006623, effort.Id, "incorrect effort id")
	assert.Equal(t, 360, effort.MovingTime, "incorrect moving time")
	assert.Equal(t, 370, effort.ElapsedTime, "incorrect elapsed time")
	assert.Equal(t, "http://pics.com/227615/large.jpg", effort.AthleteProfile, "incorrect athelete profile pic")

}
