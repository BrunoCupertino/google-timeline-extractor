package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type TimeLineObjects struct {
	TimeLineObjects []TimeLineObject `json:"timelineObjects"`
}

type Location struct {
	Name string `json:"name"`
}

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r, _ := strconv.Unquote(string(s))
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*(*time.Time)(t) = time.Unix(0, q*int64(time.Millisecond))
	return nil
}

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the JSON time as a time.Time instance in UTC
func (t Time) Time() time.Time {
	return time.Time(t).UTC()
}

// String returns t as a formatted string
func (t Time) String() string {
	return t.Time().String()
}

type Duration struct {
	StartTimestampMs Time `json:"startTimestampMs"`
	EndTimestampMs   Time `json:"endTimestampMs"`
}

type PlaceVisit struct {
	Location Location `json:"location"`
	Duration Duration `json:"duration"`
}

type TimeLineObject struct {
	ActivitySegment map[string]interface{} `json:"activitySegment"`
	PlaceVisit      PlaceVisit             `json:"placeVisit"`
}

func main() {
	jsonFile, err := os.Open("history/2017/2017_JULY.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully opened file")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result TimeLineObjects

	json.Unmarshal([]byte(byteValue), &result)

	for i := 0; i < len(result.TimeLineObjects); i++ {
		pv := result.TimeLineObjects[i].PlaceVisit
		if pv.Location.Name == os.Args[1] {
			fmt.Println("from:", pv.Duration.StartTimestampMs.Time().Local(), "to:", pv.Duration.EndTimestampMs.Time().Local(), "=====> diff:", pv.Duration.EndTimestampMs.Time().Sub(pv.Duration.StartTimestampMs.Time()))
		}
	}
}
