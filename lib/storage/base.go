package storage

import "github.com/spf13/cobra"

type StorageBackend interface {
	ConfigFlags(...*cobra.Command)
	ExportCmd() *cobra.Command
	Start(inpArgs []string)
	PrintLink()
	LoadLink()
}
