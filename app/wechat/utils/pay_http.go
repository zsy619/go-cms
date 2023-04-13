package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpServicePost(xml, url string, isUseCert bool, timeout int) ([]byte, error) {
	body := bytes.NewBuffer([]byte(xml))
	response, err := http.Post(url, "text/xml", body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
