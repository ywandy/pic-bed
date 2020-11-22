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
	cmdBase64 := &cobra.Command{
		Use:   "link [storage type]",
		Short: "get config link",
	}
	rootCmd.AddCommand(cmdSave, cmdBase64)
	cmdSave.AddCommand(s3.CmdInit())
	cmdBase64.AddCommand(s3.CmdBase64Init())
	rootCmd.Execute()
}
