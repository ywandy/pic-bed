package s3Storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"pic-bed/lib/reporter"
	"pic-bed/lib/storage"
	"time"
)

type MinioConfig struct {
	SK     string `json:"sk"`
	AK     string `json:"ak"`
	SSL    bool   `json:"ssl"`
	Host   string `json:"host"`
	Bucket string `json:"bucket"`
}

type S3Storage struct {
	Config   MinioConfig
	setting  s3Setting
	minioCli *minio.Client
}

type s3Setting struct {
	Print     bool
	QuickLink string
}

var DefaultReporter reporter.TextReporter

func (s *S3Storage) ConfigFlags(cmds ...*cobra.Command) {
	cfg := &s.Config
	setting := &s.setting
	for _, c := range cmds {
		//配置文件
		c.Flags().StringVarP(&cfg.SK, "sk", "", "", "s3 config sk")
		c.Flags().StringVarP(&cfg.AK, "ak", "", "", "s3 config ak")
		c.Flags().StringVarP(&cfg.Host, "host", "", "", "s3 config host")
		c.Flags().StringVarP(&cfg.Bucket, "bucket", "", "", "s3 config bucket")
		c.Flags().BoolVarP(&cfg.SSL, "ssl", "", false, "s3 config ssl")
		//设置
		c.Flags().BoolVarP(&setting.Print, "print", "p", false, "print config to link")
	}
}

func (s *S3Storage) ExportCmd() *cobra.Command {
	saveTypeS3 := &cobra.Command{
		Use:     "s3 [pic paths]",
		Short:   "s3 storage backend",
		Args:    cobra.MinimumNArgs(1),
		Example: "xxx.jpg -sk val1 -ak val2 -host val3 -bucket val4 -ssl false",
		Run: func(cmd *cobra.Command, args []string) {
			s.Start(args)
		},
	}
	s.ConfigFlags(saveTypeS3)
	return saveTypeS3
}

func (s *S3Storage) Start(inpArgs []string) {
	var err error
	DefaultReporter = reporter.TyporaReporter()
	cfg := s.Config
	s.minioCli, err = minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AK, cfg.SK, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		DefaultReporter(reporter.ReporterSchema{
			FileUrl: nil,
			Error:   err,
		}).Print()
		return
	}
	urlList, err := s.startUpload(inpArgs)
	if err != nil {
		DefaultReporter(reporter.ReporterSchema{
			FileUrl: nil,
			Error:   err,
		}).Print()
	} else {
		DefaultReporter(reporter.ReporterSchema{
			FileUrl: urlList,
			Error:   nil,
		}).Print()
	}
}

var year = time.Now().Year()
var month = time.Now().Month()

func (s *S3Storage) startUpload(p []string) ([]string, error) {
	fUrls := make([]string, 0)
	for _, uri := range p {
		save, err := storage.ContentFromPath(uri)
		if err != nil {
			fUrls = append(fUrls, err.Error())
			continue
		}
		ext := storage.IsKnownContentType(save.ContentType)
		fName := fmt.Sprintf("%d%d/%d", year, month, save.Timestamp.UnixNano())
		if ext != "" {
			fName = fName + "." + ext
		}
		//开始上传
		if _, err := s.minioCli.PutObject(context.Background(), s.Config.Bucket, fName, bytes.NewReader(save.Data), -1, minio.PutObjectOptions{
			ContentType: save.ContentType,
		}); err != nil {
			fUrls = append(fUrls, "error:"+err.Error()+"\n")
		} else {
			schema := "http"
			if s.Config.SSL {
				schema = "https"
			}
			fUrls = append(fUrls, schema+"://"+s.Config.Host+"/"+s.Config.Bucket+"/"+fName+"\n")
		}
	}
	return fUrls, nil
}
