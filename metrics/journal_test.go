package metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestDiskMetric(t *testing.T) {
	t.Run("should collect all logs until one minute from now", func(t *testing.T) {
		mockLogs := []string{"log1", "log2", "log3"}
		mockJournalCollector := &MockJournalCollector{MockData: mockLogs}
		journalLogHandler := &MockJournalLogsHandler{}

		previousEnd := getTime("2024-08-12T08:19:00Z")
		currentTime := getTime("2024-08-12T08:22:30Z")

		collectJournalStruct := &CollectJournalStruct{
			PreviousEndTime:    previousEnd,
			CurrentTime:        currentTime,
			CollectJournalLogs: mockJournalCollector.CollectJournalLogs,
			HandleJournalLogs:  journalLogHandler.HandleJournalLogs,
		}

		expectedEndTime := getTime("2024-08-12T08:22:00Z")
		endTime := CollectJournalLogsAndReturnLastEndTime(collectJournalStruct)
		assert.Equal(t, 3, len(journalLogHandler.Logs))
		assert.EqualValues(t, mockLogs, journalLogHandler.Logs)
		mockJournalCollector.PrintStatEndTimes()
		assert.EqualValues(t, []MockStartEndTime{
			{Start: getTime("2024-08-12T08:19:00Z"), End: getTime("2024-08-12T08:20:00Z")},
			{Start: getTime("2024-08-12T08:20:00Z"), End: getTime("2024-08-12T08:21:00Z")},
			{Start: getTime("2024-08-12T08:21:00Z"), End: getTime("2024-08-12T08:22:00Z")},
		}, mockJournalCollector.StartEndTimes)
		assert.Equal(t, expectedEndTime, endTime)
		assert.Equal(t, 3, len(mockJournalCollector.StartEndTimes))
	})

	t.Run("previous end time is 0", func(t *testing.T) {
		mockLogs := []string{"log1"}
		mockJournalCollector := &MockJournalCollector{MockData: mockLogs}
		journalLogHandler := &MockJournalLogsHandler{}

		previousEnd := time.Time{}
		currentTime := getTime("2024-08-12T08:22:30Z")

		collectJournalStruct := &CollectJournalStruct{
			PreviousEndTime:    previousEnd,
			CurrentTime:        currentTime,
			CollectJournalLogs: mockJournalCollector.CollectJournalLogs,
			HandleJournalLogs:  journalLogHandler.HandleJournalLogs,
		}

		expectedEndTime := getTime("2024-08-12T08:22:00Z")
		endTime := CollectJournalLogsAndReturnLastEndTime(collectJournalStruct)
		assert.Equal(t, 1, len(journalLogHandler.Logs))
		assert.EqualValues(t, mockLogs, journalLogHandler.Logs)
		mockJournalCollector.PrintStatEndTimes()
		assert.Equal(t, expectedEndTime, endTime)
		assert.Equal(t, 1, len(mockJournalCollector.StartEndTimes))
		assert.EqualValues(t, []MockStartEndTime{
			{Start: getTime("2024-08-12T08:21:00Z"), End: getTime("2024-08-12T08:22:00Z")},
		}, mockJournalCollector.StartEndTimes)
	})
}

func getTime(value string) time.Time {
	timeObject, err := time.Parse("2006-01-02T15:04:05Z", value)
	checkError(err)
	return timeObject
}

func mockGetTimeMinutesAgoOnTheFullMinute(minutes int) time.Time {
	nowTime, err := time.Parse("2006-01-02T15:04:05Z", "2024-08-12T10:00:00Z")
	checkError(err)
	return nowTime.Add(time.Duration(-minutes) * time.Minute).Truncate(time.Minute)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

type MockStartEndTime struct {
	Start time.Time
	End   time.Time
}

type MockJournalCollector struct {
	CallIndex     int
	MockData      []string
	StartEndTimes []MockStartEndTime
}

func (m *MockJournalCollector) CollectJournalLogs(start time.Time, end time.Time) string {
	m.StartEndTimes = append(m.StartEndTimes, MockStartEndTime{Start: start, End: end})
	data := m.MockData[m.CallIndex]
	m.CallIndex++
	return data
}

func (m *MockJournalCollector) PrintStatEndTimes() {
	println("MockJournalCollector StartEndTimes")
	for _, timeObject := range m.StartEndTimes {
		fmt.Printf("Start: %v, End %v\n", timeObject.Start.String(), timeObject.End.String())
	}
}

type MockJournalLogsHandler struct {
	Logs []string
}

func (m *MockJournalLogsHandler) HandleJournalLogs(logs string) {
	m.Logs = append(m.Logs, logs)
}
