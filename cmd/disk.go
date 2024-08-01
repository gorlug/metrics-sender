package cmd

import (
	"github.com/spf13/cobra"
	"metrics-sender/metrics"
)

const cmdName = "fileSystems"

var fileSystems []string
var isZpool bool

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "send a disk metric",
	Run: func(cmd *cobra.Command, args []string) {
		metrics.SendDiskMetrics(fileSystems, isZpool, GetConfigUrl())
	},
}

func init() {
	RootCmd.AddCommand(diskCmd)

	diskCmd.Flags().StringSliceVar(&fileSystems, cmdName, []string{}, "file systems to check")
	err := diskCmd.MarkFlagRequired(cmdName)
	diskCmd.Flags().BoolVar(&isZpool, "zpool", false, "check zpool instead of file systems")
	if err != nil {
		panic(err)
	}
}
