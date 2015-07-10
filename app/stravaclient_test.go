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

func Test_GetActivitySegments(t *testing.T) {

	mockHttpClient := MockHttpClient{responseBody: []byte(testActivity)}

	setHttpClient(&mockHttpClient)

	activity, _ := GetActivity(1)

	assert.Equal(t, 321934, activity.Id, "incorrect activity id")
	assert.Equal(t, 1, len(activity.SegmentEfforts), "incorrect segment effort length")
	assert.Equal(t, 543755075, activity.SegmentEfforts[0].Id, "incorrect segment effort id")
	assert.Equal(t, Segment{2417854, "Dash for the Ferry", 1055.11, []float32{37.79536649, -122.2796434}, []float32{37.7905785,
		-122.27015622}}, activity.SegmentEfforts[0].Segment, "incorrect segment details")

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

var testActivity = `
{
  "id": 321934,
  "resource_state": 3,
  "external_id": "2012-12-12_21-40-32-80-29011.fit",
  "upload_id": 361720,
  "athlete": {
    "id": 227615,
    "resource_state": 1
  },
  "name": "12/12/2012 San Francisco",
  "description": "the best ride ever",
  "distance": 4475.4,
  "moving_time": 1303,
  "elapsed_time": 1333,
  "total_elevation_gain": 8.6,
  "type": "Run",
  "start_date": "2012-12-13T03:43:19Z",
  "start_date_local": "2012-12-12T19:43:19Z",
  "timezone": "(GMT-08:00) America/Los_Angeles",
  "start_latlng": [
    37.8,
    -122.27
  ],
  "end_latlng": [
    37.8,
    -122.27
  ],
  "location_city": "San Francisco",
  "location_state": "CA",
  "location_country": "United States",
  "achievement_count": 6,
  "kudos_count": 1,
  "comment_count": 1,
  "athlete_count": 1,
  "photo_count": 0,
  "total_photo_count": 0,
  "photos": {
    "count": 2,
    "primary": {
      "id": null,
      "source": 1,
      "unique_id": "d64643ec9205",
      "urls": {
        "100": "http://pics.com/28b9d28f-128x71.jpg",
        "600": "http://pics.com/28b9d28f-768x431.jpg"
      }
    }
  },
  "map": {
    "id": "a32193479",
    "polyline": "kiteFpCBCD]",
    "summary_polyline": "{cteFjcaBkCx@gEz@",
    "resource_state": 3
  },
  "trainer": false,
  "commute": false,
  "manual": false,
  "private": false,
  "flagged": false,
  "workout_type": 2,
  "gear": {
    "id": "g138727",
    "primary": true,
    "name": "Nike Air",
    "distance": 88983.1,
    "resource_state": 2
  },
  "average_speed": 3.4,
  "max_speed": 4.514,
  "calories": 390.5,
  "has_kudoed": false,
  "segment_efforts": [
    {
      "id": 543755075,
      "resource_state": 2,
      "name": "Dash for the Ferry",
      "segment": {
        "id": 2417854,
        "resource_state": 2,
        "name": "Dash for the Ferry",
        "activity_type": "Run",
        "distance": 1055.11,
        "average_grade": -0.1,
        "maximum_grade": 2.7,
        "elevation_high": 4.7,
        "elevation_low": 2.7,
        "start_latlng": [
          37.7905785,
          -122.27015622
        ],
        "end_latlng": [
          37.79536649,
          -122.2796434
        ],
        "climb_category": 0,
        "city": "Oakland",
        "state": "CA",
        "country": "United States",
        "private": false
      },
      "activity": {
        "id": 32193479,
        "resource_state": 1
      },
      "athlete": {
        "id": 3776,
        "resource_state": 1
      },
      "kom_rank": 2,
      "pr_rank": 1,
      "elapsed_time": 304,
      "moving_time": 304,
      "start_date": "2012-12-13T03:48:14Z",
      "start_date_local": "2012-12-12T19:48:14Z",
      "distance": 1052.33,
      "start_index": 5348,
      "end_index": 6485,
      "hidden": false,
      "achievements": [
        {
          "type_id": 2,
          "type": "overall",
          "rank": 2
        },
        {
          "type_id": 3,
          "type": "pr",
          "rank": 1
        }
      ]
    }
  ],
  "splits_metric": [
    {
      "distance": 1002.5,
      "elapsed_time": 276,
      "elevation_difference": 0,
      "moving_time": 276,
      "split": 1
    },
    {
      "distance": 475.7,
      "elapsed_time": 139,
      "elevation_difference": 0,
      "moving_time": 139,
      "split": 5
    }
  ],
  "splits_standard": [
    {
      "distance": 1255.9,
      "elapsed_time": 382,
      "elevation_difference": 3.2,
      "moving_time": 382,
      "split": 3
    }
  ],
  "best_efforts": [
    {
      "id": 273063933,
      "resource_state": 2,
      "name": "400m",
      "segment": null,
      "activity": {
        "id": 32193479
      },
      "athlete": {
        "id": 3776
      },
      "kom_rank": null,
      "pr_rank": null,
      "elapsed_time": 105,
      "moving_time": 106,
      "start_date": "2012-12-13T03:43:19Z",
      "start_date_local": "2012-12-12T19:43:19Z",
      "distance": 400,
      "achievements": [

      ]
    },
    {
      "id": 273063935,
      "resource_state": 2,
      "name": "1/2 mile",
      "segment": null,
      "activity": {
        "id": 32193479
      },
      "athlete": {
        "id": 3776
      },
      "kom_rank": null,
      "pr_rank": null,
      "elapsed_time": 219,
      "moving_time": 220,
      "start_date": "2012-12-13T03:43:19Z",
      "start_date_local": "2012-12-12T19:43:19Z",
      "distance": 805,
      "achievements": [

      ]
    }
  ]
}`
