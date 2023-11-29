package mp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type AccessToken struct {
	AppId     string
	AppSecret string
	TmpName   string
	LckName   string
}

// get fresh access_token string
func (tkn *AccessToken) Fresh() (string, error) {
	if tkn.TmpName == "" {
		tkn.TmpName = tkn.AppId + "-accesstoken.tmp"
	}
	if tkn.LckName == "" {
		tkn.LckName = tkn.TmpName + ".lck"
	}
	for {
		if tkn.locked() {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	fi, err := os.Stat(tkn.TmpName)
	if err != nil && !os.IsExist(err) {
		return tkn.fetchAndStore()
	}
	expires := fi.ModTime().Add(2 * time.Hour).Unix()
	if expires <= time.Now().Unix() {
		return tkn.fetchAndStore()
	}
	tmp, err := os.Open(tkn.TmpName)
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	data, err := io.ReadAll(tmp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (tkn *AccessToken) fetchAndStore() (string, error) {
	if err := tkn.lock(); err != nil {
		return "", err
	}
	defer tkn.unlock()
	token, err := tkn.fetch()
	if err != nil {
		return "", err
	}
	logs.Debug(token)
	if err := tkn.store(token); err != nil {
		return "", err
	}
	return token, nil
}

func (tkn *AccessToken) store(token string) error {
	path := path.Dir(tkn.TmpName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	tmp, err := os.OpenFile(tkn.TmpName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer tmp.Close()
	if _, err := tmp.Write([]byte(token)); err != nil {
		return err
	}
	return nil
}

func (tkn *AccessToken) fetch() (string, error) {
	rtn, err := get(fmt.Sprintf("%stoken?grant_type=client_credential&appid=%s&secret=%s",
		UrlPrefix,
		tkn.AppId,
		tkn.AppSecret,
	))
	if err != nil {
		return "", err
	}
	return rtn.AccessToken, nil
}

func (tkn *AccessToken) unlock() error {
	return os.Remove(tkn.LckName)
}

func (tkn *AccessToken) lock() error {
	path := path.Dir(tkn.LckName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	lck, err := os.Create(tkn.LckName)
	if err != nil {
		return err
	}
	lck.Close()
	return nil
}

func (tkn *AccessToken) locked() bool {
	_, err := os.Stat(tkn.LckName)
	return !os.IsNotExist(err)
}
