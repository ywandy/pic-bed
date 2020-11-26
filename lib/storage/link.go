package storage

import (
	"encoding/base64"
	"github.com/vmihailenco/msgpack/v5"
)

func MarshToMsgPackString(i interface{}) string {
	n, _ := msgpack.Marshal(i)
	sEnc := base64.StdEncoding.EncodeToString(n)
	return sEnc
}

func UnMarshMsgPackStringToStuct(str string, i interface{}) error {
	sDec, _ := base64.StdEncoding.DecodeString(str)
	return msgpack.Unmarshal(sDec, &i)
}
