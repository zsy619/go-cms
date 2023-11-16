package mp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type res struct {
	// error fields
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// response from weixinmp
type response struct {
	// error fields
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// token fields
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	// media fields
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	// ticket fields
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

func post(url string, bodyType string, body *bytes.Buffer) (*response, error) {
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, fmt.Errorf("%d %s", rtn.ErrCode, rtn.ErrMsg)
	}
	return &rtn, nil
}

func get(url string) (*response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, fmt.Errorf("%d %s", rtn.ErrCode, rtn.ErrMsg)
	}
	return &rtn, nil
}

func postjson(surl, jsonstr string) (b []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", surl, bytes.NewReader([]byte(jsonstr)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()
	b, err = io.ReadAll(r.Body)
	return
}

func getbytes(surl string) (b []byte, err error) {
	resp, err := http.Get(surl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	return
}
