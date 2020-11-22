package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var minioCli *minio.Client

var result = "Upload Success:\n"

func Start(cfg MinioConfig, inpArgs []string) {

	var err error
	minioCli, err = minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AK, cfg.SK, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		result += err.Error() + "\n"
		fmt.Printf(result)
		return
	}
	Uploader(cfg, inpArgs)
}

var year = time.Now().Year()
var month = time.Now().Month()
var date = time.Now().Format("20060102")
var imgPerfix = []string{"png", "jpg", "jpeg", "gif", "bmp", "tif", "webp", "exif"}

func download(imgurl string) ([]byte, string, error) {
	//判断url合法
	rsp, err := http.Get(imgurl)
	if err != nil {
		return []byte{}, "", err
	}
	defer rsp.Body.Close()
	//获取content-type
	perfixOk := false
	perfix := ""
	contentType, ok := rsp.Header["Content-Type"]
	if ok && len(contentType) > 0 {
		for _, p := range imgPerfix {
			if strings.HasSuffix(contentType[0], fmt.Sprintf("/%s", p)) {
				perfixOk = true
				perfix = p
				continue
			}
		}
	}
	if perfixOk {
		//如果后缀合适的情况下
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return []byte{}, "", err
		}
		return b, perfix, nil
	}
	return []byte{}, "", errors.New("Content-Type not a Image")
}

func getContentType(ct string) string {
	var base = "image/"
	if ct == "" {
		return base + "jpg"
	} else {
		return base + ct
	}
}

func Uploader(cfg MinioConfig, inpArgs []string) {
	for _, img := range inpArgs {
		u, err := url.Parse(img)
		timestamp := time.Now().UnixNano()
		fName := fmt.Sprintf("%d%d/%s_%d", year, month, date, timestamp)
		if err != nil {
			result += "error:" + err.Error() + "\n"
			continue
		}
		if strings.HasPrefix(u.Scheme, "http") {
			b, fPerfix, err := download(u.String())
			fName = fmt.Sprintf("%s.%s", fName, fPerfix)
			if err != nil {
				result += "error:" + err.Error() + "\n"
			} else {
				if _, err := minioCli.PutObject(context.Background(), cfg.Bucket, fName, bytes.NewReader(b), -1, minio.PutObjectOptions{
					ContentType: getContentType(fPerfix),
				}); err != nil {
					result += "error:" + err.Error() + "\n"
				} else {
					result += "https://" + cfg.Host + "/" + cfg.Bucket + "/" + fName + "\n"
				}
			}
		} else {
			fPerfix := strings.ReplaceAll(path.Ext(img), ".", "")
			fName = fmt.Sprintf("%s.%s", fName, fPerfix)
			if _, err := minioCli.FPutObject(context.Background(), cfg.Bucket, fName, img, minio.PutObjectOptions{
				ContentType: getContentType(fPerfix),
			}); err != nil {
				result += "error:" + err.Error() + "\n"
			} else {
				result += "https://" + cfg.Host + "/" + cfg.Bucket + "/" + fName + "\n"
			}
		}
	}
	fmt.Println(result)
}
