package s3

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var minioConfig *MinioConfig

func CmdInit() *cobra.Command {
	var picUrl = ""
	saveTypeS3 := &cobra.Command{
		Use:   "s3 [pic paths]",
		Short: "s3 storage backend",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Flags().String("url", "", "base64 url"))
			if picUrl != "" {
				json.Unmarshal([]byte(DecodeBase64Url(picUrl)), minioConfig)
			}
			Start(*minioConfig, args)
		},
	}
	saveTypeS3.Flags().String("url", "", "base64 url")
	minioConfig = CmdSetMinioConfigFlags(saveTypeS3)
	return saveTypeS3
}

func DecodeHook() {

}

func DecodeBase64Url(s string) string {
	out := "{}"
	sDec, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(s, "pic://", ""))
	if err == nil {
		out = string(sDec)
	}
	return out
}

func CmdBase64Init() *cobra.Command {
	printBase64S3 := &cobra.Command{
		Use:   "s3 [pic paths]",
		Short: "s3 storage backend",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := json.Marshal(minioConfig)
			uEnc := base64.URLEncoding.EncodeToString(b)
			fmt.Println(fmt.Sprintf("pic://%s", uEnc))
		},
	}
	minioConfig = CmdSetMinioConfigFlags(printBase64S3)
	return printBase64S3
}
