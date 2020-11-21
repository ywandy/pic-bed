package s3

import (
	"github.com/spf13/cobra"
)

func CmdInit() *cobra.Command {
	var sk string
	var ak string
	var host string
	var bucket string
	var ssl bool
	saveTypeS3 := &cobra.Command{
		Use:   "s3 [pic paths]",
		Short: "s3 storage backend",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Start(minioConfig{
				Sk:     sk,
				Ak:     ak,
				Host:   host,
				Bucket: bucket,
				Ssl:    ssl,
			}, args)
		},
	}
	saveTypeS3.Flags().StringVarP(&sk, "sk", "", "", "s3 config sk")
	saveTypeS3.Flags().StringVarP(&ak, "ak", "", "", "s3 config ak")
	saveTypeS3.Flags().StringVarP(&host, "host", "", "", "s3 config host")
	saveTypeS3.Flags().StringVarP(&bucket, "bucket", "", "", "s3 config bucket")
	saveTypeS3.Flags().BoolVarP(&ssl, "ssl", "", false, "s3 config ssl")
	return saveTypeS3
}
