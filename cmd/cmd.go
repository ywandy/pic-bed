package cmd

import (
	"github.com/spf13/cobra"
	"pic-bed/lib/s3"
)

var ConfigPath = "./config.yaml"

func init() {
	var rootCmd = &cobra.Command{Use: "picbed"}
	cmdSave := &cobra.Command{
		Use:   "save [storage type]",
		Short: "save to a storage backend",
	}
	rootCmd.AddCommand(cmdSave)
	cmdSave.AddCommand(s3.CmdInit())
	rootCmd.Execute()
}
