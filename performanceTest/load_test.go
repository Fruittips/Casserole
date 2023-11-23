package performanceTests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type requestData struct {
	StudentName string `json:"studentName"`
	StudentId   string `json:"studentId"`
	CourseId    string `json:"courseId"`
}

var portsStr = flag.String("ports", "", "Ports to hit")

func TestPerformance(t *testing.T) {
	flag.Parse()
	ports := strings.Split(*portsStr, ",")
	rounds := 10
	record := make([]map[string]time.Duration, rounds)

	for round := 0; round < rounds; round++ {
		var wg sync.WaitGroup

		// round 1 = 1 write + 1 read, round 2 = 2 write + 2 read, etc.
		numberOfReadsAndWrites := round + 1
		wg.Add(numberOfReadsAndWrites * 2)

		// channels to store timings
		writeTimings := make(chan time.Duration, numberOfReadsAndWrites)
		readTimings := make(chan time.Duration, numberOfReadsAndWrites)

		// we assume that we have a load balancer that hits nodes in a clockwise round robin manner
		for r := 0; r < numberOfReadsAndWrites; r++ {
			go func(p string) {
				defer wg.Done()
				startTime := time.Now()
				executeReadReq(p)
				readTimings <- time.Since(startTime)
			}(ports[r%len(ports)])

			go func(p string) {
				defer wg.Done()
				startTime := time.Now()
				executeWriteReq(p)
				writeTimings <- time.Since(startTime)
			}(ports[r%len(ports)])
		}

		wg.Wait()
		close(writeTimings)
		close(readTimings)

		var totalPostTime, totalGetTime time.Duration
		for t := range writeTimings {
			totalPostTime += t
		}
		for t := range readTimings {
			totalGetTime += t
		}

		// calculate average time, time.Duration(numberOfWrites) is in nanosecond which is ok as totalPostTime is also in nanosecond
		averagePostTime := totalPostTime / time.Duration(numberOfReadsAndWrites)
		averageGetTime := totalGetTime / time.Duration(numberOfReadsAndWrites)
		record[round] = map[string]time.Duration{
			"Writes": averagePostTime,
			"Reads":  averageGetTime,
		}
	}

	for i, r := range record {
		fmt.Printf("Round %d: %v\n", i+1, r)
	}
}

func executeWriteReq(port string) (resp *http.Response, err error) {
	data := requestData{
		StudentName: "Mah Yi Da",
		StudentId:   "123123",
		CourseId:    "CS-101",
	}
	jsonData, _ := json.Marshal(data)
	url := fmt.Sprintf("http://127.0.0.1:%s/write/course/%s", port, data.CourseId)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func executeReadReq(port string) (resp *http.Response, err error) {
	data := requestData{
		StudentName: "Mah Yi Da",
		StudentId:   "123123",
		CourseId:    "CS-101",
	}

	url := fmt.Sprintf("http://127.0.0.1:%s/read/course/%s/student/%s", port, data.CourseId, data.StudentId)
	return http.Get(url)
}
