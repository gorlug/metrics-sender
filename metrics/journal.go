package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"metrics-backend/rest"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func SendJournalLogs(metaInfoPath string, url string) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		println(fmt.Sprintf("failed to load location %v, error: %v", "Europe/Berlin", err))
		location = time.Local
	}

	previousEndTime := time.Time{}
	metaInfo, err := os.ReadFile(metaInfoPath)
	if err == nil {
		parsedTime, err := parseTime(metaInfo)
		handleError(err)
		previousEndTime = parsedTime.In(location).Truncate(time.Minute)
	}
	collectJournalStruct := &CollectJournalStruct{
		PreviousEndTime:    previousEndTime,
		CurrentTime:        time.Now().In(location).Truncate(time.Minute),
		CollectJournalLogs: CollectJournalLogs,
		HandleJournalLogs: func(logs string) {
			HandleJournalLogs(logs, url)
		},
	}
	lastEndTime := CollectJournalLogsAndReturnLastEndTime(collectJournalStruct)
	err = os.WriteFile(metaInfoPath, []byte(formatTime(lastEndTime.UTC())), 0644)
}

func parseTime(metaInfo []byte) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(metaInfo))
}

type CollectJournalStruct struct {
	PreviousEndTime    time.Time
	CurrentTime        time.Time
	CollectJournalLogs func(start time.Time, end time.Time) string
	HandleJournalLogs  func(logs string)
}

func CollectJournalLogsAndReturnLastEndTime(values *CollectJournalStruct) time.Time {
	maxEndTime := values.CurrentTime.Truncate(time.Minute)
	endTime := values.PreviousEndTime
	if endTime == (time.Time{}) {
		endTime = maxEndTime.Add(time.Duration(-1) * time.Minute)
	}

	var startTime time.Time
	for {
		startTime = endTime
		endTime = startTime.Add(time.Duration(1) * time.Minute)
		println(fmt.Sprintf("Before check Start: %v, End: %v\n, maxEndTime: %v", startTime, endTime, maxEndTime))
		if endTime.After(maxEndTime) {
			break
		}
		println(fmt.Sprintf("Start: %v, End: %v\n, maxEndTime: %v", startTime, endTime, maxEndTime))
		logs := values.CollectJournalLogs(startTime, endTime)
		values.HandleJournalLogs(logs)
	}
	return maxEndTime
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func CollectJournalLogs(start time.Time, end time.Time) string {
	out, err := exec.Command("journalctl",
		"--since", formatTime(start),
		"--until", formatTime(end), "-o", "json").Output()
	handleError(err)
	return string(out)
}

func HandleJournalLogs(logs string, url string) {
	journalBody := &rest.JournalBody{
		Logs: logs,
	}
	bodyJson, marshalErr := json.Marshal(journalBody)

	if marshalErr != nil {
		log.Fatal("Could not marshal metric into JSON")
	}

	// Make request with marshalled JSON as the POST body
	response, err := http.Post(url, "application/json",
		bytes.NewBuffer(bodyJson))

	if err != nil {
		log.Fatal("Could not make POST request")
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Error Response status: %v\n", response.Status)
	}
}
