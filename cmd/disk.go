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
	Long:  "Send a a disk metric to the metrics-backend. It sends the disk usage in percent for the specified file systems.",
	Run: func(cmd *cobra.Command, args []string) {
		metrics.SendDiskMetrics(fileSystems, isZpool, GetConfigUrl())
	},
}

func init() {
	RootCmd.AddCommand(diskCmd)

	diskCmd.Flags().StringSliceVar(&fileSystems, cmdName, []string{}, "File systems to check. Provide multiple separated by a comma e.g. \"/,/mnt/storage\"")
	err := diskCmd.MarkFlagRequired(cmdName)
	diskCmd.Flags().BoolVar(&isZpool, "zpool", false, "this means the provided file systems are ZFS datasets instead of normal file systems")
	if err != nil {
		panic(err)
	}
}
