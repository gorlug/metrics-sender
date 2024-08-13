package metrics

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func SendJournalLogs(metaInfoPath string) {
	previousEndTime := time.Time{}
	metaInfo, err := os.ReadFile(metaInfoPath)
	if err == nil {
		parsedTime, err := parseTime(metaInfo)
		handleError(err)
		previousEndTime = parsedTime
	}
	collectJournalStruct := &CollectJournalStruct{
		PreviousEndTime:    previousEndTime,
		CurrentTime:        time.Now(),
		CollectJournalLogs: CollectJournalLogs,
		HandleJournalLogs:  HandleJournalLogs,
	}
	lastEndTime := CollectJournalLogsAndReturnLastEndTime(collectJournalStruct)
	err = os.WriteFile(metaInfoPath, []byte(formatTime(lastEndTime)), 0644)
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
		if endTime.After(maxEndTime) {
			break
		}
		fmt.Printf("Start: %v, End: %v\n", startTime, endTime)
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

func HandleJournalLogs(logs string) {
	fmt.Printf("Logs: %v", logs)
}
