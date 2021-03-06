package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func main() {
	file := flag.String("file", "", "path to the file to be read")
	name := flag.String("name", "*", "name of the location to be used to filter")
	output := flag.String("output", "console", "output format (console / file)")

	flag.Parse()

	fileValue := *file
	namesValue := strings.Split(strings.ToLower(*name), "|")
	outputValue := *output

	jsonFile, err := os.Open(fileValue)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully opened file")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result TimeLineObjects

	json.Unmarshal([]byte(byteValue), &result)

	outputData := [][]string{
		{"Location", "Start From", "End To", "File"},
	}

	for i := 0; i < len(result.TimeLineObjects); i++ {
		pv := result.TimeLineObjects[i].PlaceVisit
		if namesValue[0] == "*" || contains(namesValue, strings.ToLower(pv.Location.Name)) {
			values := []string{pv.Location.Name, pv.Duration.StartTimestampMs.Time().Local().Format("2006-01-02T15:04:05"), pv.Duration.EndTimestampMs.Time().Local().Format("2006-01-02T15:04:05"), fileValue}
			outputData = append(outputData, values)
		}
	}

	csvFile, err := os.OpenFile("result.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	csvWriter := csv.NewWriter(csvFile)

	for _, entry := range outputData {
		if outputValue == "file" {
			_ = csvWriter.Write(entry)
		} else {
			fmt.Printf("location: %s, start from:%s, end to:%s, file: %s\n", entry[0], entry[1], entry[2], entry[3])
		}
	}

	csvWriter.Flush()
	csvFile.Close()
}
