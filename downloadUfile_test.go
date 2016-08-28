package ufile

import (
	"net/http"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	bucket := "your_bucket_name"
	key := "the_filename_to_download"
	localfile := "the_name_of_download_file_local"

	ufile := UFile{}
	url := ufile.DownloadFile(bucket, key, localfile)

	resp, _ := http.Head(url)

	if resp.StatusCode != 200 {
		t.Error("Get Download url failed!")
	}
}
