package app

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
)

type Area struct {
	id          int
	name        string
	boundingBox BoundingBox
}

type Rider struct {
	Id   int      `json:"id"`
	Name string   `json:"name"`
	KOMs []Effort `json:"koms"`
}

var DBFileName string = "kombagger.db"

var applicationData StravaData

// StravaData is a thin wrapper around the bolt DB to provide persistence.
type StravaData interface {

	// AddKOM defines and stores this Effort as a KOM of this rider.
	AddKOM(effort Effort)

	// GetRiderLeaderboard returns a slice of Riders. The order is determined by the
	// number of KOMs held by the rider.
	GetRiderLeaderboard() []Rider

	// Delete all data from the DB
	Clean()
}

const RIDERS_BUCKET string = "Riders"

type BoltStravaData struct {
	db *bolt.DB
}

func (boltStravaData BoltStravaData) AddKOM(effort Effort) {
	var rider *Rider
	boltStravaData.db.View(func(tx *bolt.Tx) error {
		riders := tx.Bucket([]byte(RIDERS_BUCKET))
		bytes := riders.Get([]byte(strconv.Itoa(effort.AthleteId)))
		if bytes != nil {
			logger.Println("Athlete already has a KOM...")
			fmt.Println(deserialize(bytes))
			temp := deserialize(bytes)
			rider = &temp

		}
		return nil
	})

	boltStravaData.db.Update(func(tx *bolt.Tx) error {
		// create the rider if they don't exist in the DB
		if rider == nil {
			rider = &Rider{effort.AthleteId, effort.AthleteName, nil}
		}

		rider.KOMs = append(rider.KOMs, effort)

		// store the new Rider overwriting any existing
		riders := tx.Bucket([]byte(RIDERS_BUCKET))
		riders.Put([]byte(strconv.Itoa(effort.AthleteId)), serialize(rider))

		logger.Printf("added new KOM for rider %s %s", strconv.Itoa(effort.AthleteId), effort.AthleteName)
		return nil
	})
}

func (boltStravaData BoltStravaData) GetRiderLeaderboard() []Rider {
	Riders := []Rider{}
	boltStravaData.db.View(func(tx *bolt.Tx) error {
		riders := tx.Bucket([]byte(RIDERS_BUCKET))

		c := riders.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			nextRider := deserialize(v)
			Riders = append(Riders, nextRider)
		}

		return nil
	})
	return Riders

}

func (boltStravaData BoltStravaData) Clean() {
	boltStravaData.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(RIDERS_BUCKET))
		return nil
	})

	boltStravaData.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(RIDERS_BUCKET))
		return nil
	})
}

func deserialize(riderBytes []byte) Rider {
	rider_ := Rider{}
	buf := bytes.NewBuffer(riderBytes)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&rider_)
	if err != nil {
		panic(err)
	}
	return rider_
}

func serialize(rider *Rider) []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(rider)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func init() {
	fmt.Println("initializing Bolt DB")
	db, _ := bolt.Open(DBFileName, 0600, nil)

	applicationData = BoltStravaData{db}
	applicationData.Clean()
}

func AddKOM(effort Effort) {
	applicationData.AddKOM(effort)
}

func GetRiderLeaderboard(areaId int) []Rider {
	return applicationData.GetRiderLeaderboard()
}
