package metrics

import (
	"fmt"
	"log"
	"metrics-backend/metrics"
	"metrics-sender/send"
	"os/exec"
	"strings"
)

func SendDockerMetrics(url string, prefix *[]string) {
	out, err := exec.Command("docker", "ps", "--filter", "status=running", "--format", "{{.Names}}").Output()

	if err != nil {
		log.Println(err)
		panic(err)
	}

	containers := strings.Split(string(out), "\n")

	for _, container := range containers {
		if container != "" && !isContainerExcluded(container, prefix) {
			metricBuilder := metrics.NewMetricBuilder().
				WithName(fmt.Sprintf("container %v", container)).
				WithType(metrics.Ping)
			send.SendMetric(metricBuilder, url)
		}
	}
}

// func that returns false if the container name is inside the prefix array
func isContainerExcluded(container string, prefix *[]string) bool {
	for _, p := range *prefix {
		if strings.HasPrefix(container, p) {
			return true
		}
	}
	return false
}
