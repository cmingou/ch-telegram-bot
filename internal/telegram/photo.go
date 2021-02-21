package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FileStruct struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int    `json:"file_size"`
		FilePath     string `json:"file_path"`
	} `json:"result"`
}

func GetPhotoUrl(token, fileId string) (string, error) {
	fs := &FileStruct{}

	url := fmt.Sprintf("https://api.telegram.org/bot%v/getFile?file_id=%v", token, fileId)
	rsp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	if err := json.NewDecoder(rsp.Body).Decode(fs); err != nil {
		return "", err
	}

	photoLink := fmt.Sprintf("https://api.telegram.org/file/bot%v/%v", token, fs.Result.FilePath)

	return photoLink, nil
}
