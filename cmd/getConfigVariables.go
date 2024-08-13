package cmd

import "github.com/spf13/viper"

func GetConfigUrl() string {
	url := viper.Get("url")
	return url.(string)
}

func GetJournalLogMetaFile() string {
	metaFile := viper.Get("journalLogMetaFile")
	return metaFile.(string)
}

func GetJournalUrl() string {
	url := viper.Get("journalUrl")
	return url.(string)
}
