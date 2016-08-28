package ufile

import (
	"crypto/hmac"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	PUBLIC_KEY     = "your_ucloud_public_key"
	PRIVATE_KEY    = "your_ucloud_private_key"
	DEFAULT_BUCKET = "bucket_name"
	USER_AGENT     = "Golang HttpRequest/0.1-beta"
	EXPIRE         = 300
	// change depend on your need
	UCLOUD_DOWNLOAD_SUFFIX = ".ufile.ucloud.cn"
)

func Signature(bucket string, key string, method string, header map[string]string) string {
	data := method + "\n"
	data += "\n" //Content-md5 null
	data += "\n" //Content-Type null
	data += header["Expires"] + "\n"
	data += "/" + bucket + "/" + key

	h := hmac.New(sha1.New, []byte(PRIVATE_KEY))
	h.Write([]byte(data))
	sEnc := b64.StdEncoding.EncodeToString(h.Sum(nil))
	return sEnc
}

type UFile struct{}

func (u *UFile) GrabFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// bucket name is file content prefix.
// key is alias file name in ufile system
func (u *UFile) DownloadFile(bucket string, key string, localfile string) string {
	var header = map[string]string{"User-Agent": USER_AGENT}
	expire := EXPIRE + time.Now().UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond))/1000
	fmt.Println(expire)
	expire_str := strconv.FormatInt(expire, 10)
	header["Expires"] = expire_str

	signature := Signature(bucket, key, "GET", header)
	fmt.Println(signature)
	query := map[string]string{
		"UCloudPublicKey": PUBLIC_KEY,
		"Expires":         expire_str,
		"Signature":       signature,
	}
	params := url.Values{}

	for k, v := range query {
		params.Add(k, v)
	}

	query_str := params.Encode()
	s := fmt.Sprintf("http://%s%s/%s?%s", bucket, UCLOUD_DOWNLOAD_SUFFIX, key, query_str)
	log.Println(s)
	return s
}
