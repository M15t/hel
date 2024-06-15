package fshare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Query FS by query string
func (s *Fshare) Query(q string) (*QueryResponse, error) {
	url := fmt.Sprintf("%s%s%s", TimFsURL, "/api/v1/string-query-search?query=", q)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(QueryResponse)
	json.Unmarshal(body, &data)

	return data, nil
}

// GetToken get token from fshare
func (s *Fshare) GetToken() (*TokenResponse, error) {
	url := FsAPIURL + "/api/user/login"
	method := "POST"

	payload := strings.NewReader(`{
	"user_email" : "` + s.cfg.UserEmail + `",
	"password":	"` + s.cfg.UserPwd + `",
	"app_key" : "` + s.cfg.AppKey + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", s.cfg.UserAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", string(body))
	}

	data := new(TokenResponse)
	json.Unmarshal(body, &data)

	return data, nil
}

// GetFilesInFolder get files in folder
func (s *Fshare) GetFilesInFolder(token, sessionID, folderURL string) ([]*FileResponse, error) {
	url := FsAPIURL + "/api/fileops/getFolderList"
	method := "POST"

	payload := strings.NewReader(`{
	"token" : "` + token + `",
	"url":	"` + folderURL + `",
	"dirOnly": 0,
	"pageIndex": -1,
	"limit": 100
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", s.cfg.UserAgent)
	req.Header.Add("Cookie", fmt.Sprintf("session_id=%s", sessionID))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", string(body))
	}

	data := make([]*FileResponse, 0)
	json.Unmarshal(body, &data)

	return data, nil
}

// GetFileInfo get file info
func (s *Fshare) GetFileInfo(token, sessionID, fileURL string) (*FileResponse, error) {
	url := FsAPIURL + "/api/fileops/get"
	method := "POST"

	payload := strings.NewReader(`
	{"url":"` + fileURL + `",
	"token":"` + token + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, &json.InvalidUnmarshalError{}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", s.cfg.UserAgent)
	req.Header.Add("Cookie", fmt.Sprintf("session_id=%s", sessionID))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", string(body))
	}

	data := new(FileResponse)
	json.Unmarshal(body, &data)

	return data, nil
}

// GetFileLocation get file location
func (s *Fshare) GetFileLocation(token, sessionID, fileURL string) (*FileDownloadResponse, error) {
	url := FsAPIURL + "/api/session/download"
	method := "POST"

	payload := strings.NewReader(`
	{"url":"` + fileURL + `",
	"password":"",
	"token":"` + token + `",
	"zipflag":0}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, &json.InvalidUnmarshalError{}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", s.cfg.UserAgent)
	req.Header.Add("Cookie", fmt.Sprintf("session_id=%s", sessionID))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := new(FileDownloadResponse)
	json.Unmarshal(body, &data)

	return data, nil
}
