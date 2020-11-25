package storage

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type SaveBody struct {
	Data        []byte
	FileName    string
	Timestamp   time.Time
	ContentType string
}

func ContentFromPath(uri string) (SaveBody, error) {
	if IsWebUrl(uri) {
		//web资源
		b, cType, err := download(uri)
		if err != nil {
			return SaveBody{}, err
		}
		return SaveBody{
			Data:        b,
			Timestamp:   time.Now(),
			ContentType: cType,
		}, nil
	} else {
		//本地资源
		t := filepath.FromSlash(uri)
		b, err := ioutil.ReadFile(t)
		if err != nil {
			return SaveBody{}, err
		}
		cType, err := GetFileContentType(b)
		if err != nil {
			return SaveBody{}, err
		}
		return SaveBody{
			Data:        b,
			Timestamp:   time.Now(),
			ContentType: cType,
		}, nil
	}
}

func download(imgurl string) ([]byte, string, error) {
	//判断url合法
	rsp, err := http.Get(imgurl)
	if err != nil {
		return []byte{}, "", err
	}
	defer rsp.Body.Close()
	//获取content-type
	//如果后缀合适的情况下
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte{}, "", err
	}
	cType, err := GetFileContentType(b)
	if err != nil {
		return []byte{}, "", err
	}
	return b, cType, nil
}
