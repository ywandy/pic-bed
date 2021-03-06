package cmd

import (
	"github.com/spf13/cobra"
	"pic-bed/lib/plugin/githubStorage"
	"pic-bed/lib/plugin/s3Storage"
	"pic-bed/lib/storage"
)

func init() {
	var rootCmd = &cobra.Command{Use: "picbed"}
	cmdSave := &cobra.Command{
		Use:   "save [storage type]",
		Short: "save to a storage backend",
	}
	rootCmd.AddCommand(cmdSave)
	s3StorageInstance := storage.StorageBackend(&s3Storage.S3Storage{})
	gitHubStorageInstance := storage.StorageBackend(&githubStorage.GitHubStorage{})
	cmdSave.AddCommand(s3StorageInstance.ExportCmd(), gitHubStorageInstance.ExportCmd())
	rootCmd.Execute()
}
