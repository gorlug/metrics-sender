package cmd

import (
	"github.com/spf13/cobra"
	"metrics-backend/metrics"
	"metrics-sender/send"
	"strconv"
)

var name string
var minutesTillAlert int

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "send a ping metric",
	Run: func(cmd *cobra.Command, args []string) {
		metricsBuilder := metrics.NewMetricBuilder().WithType(metrics.Ping).WithName(name)
		if minutesTillAlert > 0 {
			metricsBuilder.WithValue(strconv.Itoa(minutesTillAlert))
		}
		send.SendMetric(metricsBuilder, GetConfigUrl())
	},
}

func init() {
	RootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringVar(&name, "name", "ping", "name of the ping")
	pingCmd.Flags().IntVar(&minutesTillAlert, "minutesTillAlert", 0, "minutes till alert")
}
