package cmd

import (
	"github.com/spf13/cobra"
	"metrics-sender/metrics"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "send journal logs",
	Run: func(cmd *cobra.Command, args []string) {
		metrics.SendJournalLogs(GetJournalLogMetaFile(), GetJournalUrl())
	},
}

func init() {
	RootCmd.AddCommand(journalCmd)
}
