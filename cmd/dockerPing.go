package cmd

import (
	"github.com/spf13/cobra"
	"metrics-sender/metrics"
)

var excludePrefix []string

var dockerPingCmd = &cobra.Command{
	Use:   "dockerPing",
	Short: "send a docker ping",
	Long:  "send a docker ping metric for every running container",
	Run: func(cmd *cobra.Command, args []string) {
		metrics.SendDockerMetrics(GetConfigUrl(), &excludePrefix)
	},
}

func init() {
	RootCmd.AddCommand(dockerPingCmd)

	dockerPingCmd.Flags().StringSliceVar(&excludePrefix, "exclude", []string{}, "exclude containers starting with this name")
}
