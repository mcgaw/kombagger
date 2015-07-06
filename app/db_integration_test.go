package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const dBFileName = "testkombagger.db"

func init() {
	err := os.Remove(dBFileName)
	if err != nil && !os.IsNotExist(err) {
		panic("unable to remove old db file")
	}
	DBFileName = dBFileName
}

func TestRiderPersistence(t *testing.T) {

	applicationData.Clean()

	testRider := Rider{1, "Strava Rider", nil}

	testRider.KOMs = append(testRider.KOMs, Effort{Id: 1, AthleteName: "Strava Rider", AthleteId: 1,
		ElapsedTime: 350, MovingTime: 350, Rank: 1, AthleteProfile: "http://pic"})

	applicationData.AddKOM(testRider.KOMs[0])

	leaderboard := applicationData.GetRiderLeaderboard()

	assert.Equal(t, 1, len(leaderboard))

	rider := leaderboard[0]
	assert.Equal(t, 1, rider.Id, "incorrect rider id")
	assert.Equal(t, "Strava Rider", rider.Name, "incorrect rider name")
	assert.Equal(t, 1, len(rider.KOMs), "incorrect number of KOMs")
	kom := rider.KOMs[0]
	assert.Equal(t, 1, kom.Id, "effort id incorrect")
	assert.Equal(t, "Strava Rider", kom.AthleteName, "athelete name incorrect")
	assert.Equal(t, 350, kom.ElapsedTime, "elapsed time incorrect")
	assert.Equal(t, 350, kom.MovingTime, "moving time incorrect")
	assert.Equal(t, 1, kom.Rank, "rank incorrect")
	assert.Equal(t, "http://pic", kom.AthleteProfile, "athlete profile incorrect")

}

func TestRiderMultipleKOMs(t *testing.T) {

	applicationData.Clean()

	testRider := Rider{1, "Strava Rider", nil}

	kom1 := Effort{Id: 1, AthleteName: "Strava Rider", AthleteId: 1,
		ElapsedTime: 350, MovingTime: 350, Rank: 1, AthleteProfile: "http://pic"}

	kom2 := Effort{Id: 2, AthleteName: "Strava Rider", AthleteId: 1,
		ElapsedTime: 350, MovingTime: 350, Rank: 1, AthleteProfile: "http://pic"}

	testRider.KOMs = append(testRider.KOMs, kom1)
	testRider.KOMs = append(testRider.KOMs, kom2)

	applicationData.AddKOM(testRider.KOMs[0])
	applicationData.AddKOM(testRider.KOMs[1])

	leaderboard := applicationData.GetRiderLeaderboard()

	assert.Equal(t, 1, len(leaderboard))
	assert.Equal(t, 2, len(leaderboard[0].KOMs))

}
