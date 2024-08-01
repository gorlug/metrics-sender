package metrics

import (
	"fmt"
	"log"
	"metrics-backend/metrics"
	"metrics-sender/send"
	"os/exec"
	"strconv"
	"strings"
)

func SendDiskMetrics(fileSystems []string, isZfsPool bool, url string) {
	for _, fileSystem := range fileSystems {
		var usedDiskSpaceInPercent int
		if isZfsPool {
			usedDiskSpaceInPercent = GetZfsUsedDiskSpaceInPercent(fileSystem)
		} else {
			usedDiskSpaceInPercent = GetUsedDiskSpaceInPercent(fileSystem)
		}
		metricBuilder := metrics.NewMetricBuilder().WithName(fileSystem).WithType(metrics.Disk).WithValue(fmt.Sprintf("%v", usedDiskSpaceInPercent))
		send.SendMetric(metricBuilder, url)
	}
}

func GetUsedDiskSpaceInPercent(fileSystem string) int {
	out, err := exec.Command("df", "-h", fileSystem).Output()

	if err != nil {
		log.Println(err)
		return 100
	}

	output := string(out)
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return 100
	}
	usageLine := lines[1]
	fields := strings.Fields(usageLine)
	percentageField := fields[4]
	return getPercentageFromField(percentageField)
}

func getPercentageFromField(percentageField string) int {
	percentage := strings.TrimSuffix(percentageField, "%")

	percentInt, err := strconv.Atoi(percentage)
	if err != nil {
		log.Println(err)
		return 100
	}
	return percentInt
}

func GetZfsUsedDiskSpaceInPercent(poolName string) int {
	out, err := exec.Command("zpool", "list", "-H", "-o", "name,capacity", poolName).Output()

	if err != nil {
		log.Println(err)
		return 100
	}

	output := string(out)
	parts := strings.Fields(output)

	return getPercentageFromField(parts[1])
}
