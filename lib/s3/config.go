package s3

import "github.com/spf13/cobra"

type MinioConfig struct {
	SK     string `json:"sk"`
	AK     string `json:"ak"`
	SSL    bool   `json:"ssl"`
	Host   string `json:"host"`
	Bucket string `json:"bucket"`
}

func CmdSetMinioConfigFlags(cmds ...*cobra.Command) *MinioConfig {
	m := &MinioConfig{}
	for _, c := range cmds {
		c.Flags().StringVarP(&m.SK, "sk", "", "", "s3 config sk")
		c.Flags().StringVarP(&m.AK, "ak", "", "", "s3 config ak")
		c.Flags().StringVarP(&m.Host, "host", "", "", "s3 config host")
		c.Flags().StringVarP(&m.Bucket, "bucket", "", "", "s3 config bucket")
		c.Flags().BoolVarP(&m.SSL, "ssl", "", false, "s3 config ssl")
	}
	return m
}
