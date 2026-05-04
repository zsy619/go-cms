package mp

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// request from weixinmp
type Request struct {
	Token string
	// request common fields
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// message request fields
	Content      string
	MsgId        int64
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64
	Label        string
	Title        string
	Description  string
	Url          string
	Recognition  string
	// event request fields
	Event     string
	EventKey  string
	Ticket    string
	Latitude  float64
	Longitude float64
	Precision float64
}

// validate request
func (request *Request) IsValid(rw http.ResponseWriter, req *http.Request) bool {
	// if !request.checkSignature(req) {
	// 	rw.WriteHeader(http.StatusUnauthorized)
	// 	rw.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	// 	return false
	// }
	// if req.Method != "POST" {
	// 	rw.WriteHeader(http.StatusOK)
	// 	rw.Write([]byte(req.FormValue("echostr")))
	// 	return false
	// }
	if err := request.parseRequest(req); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return false
	}
	return true
}

func (request *Request) parseRequest(req *http.Request) error {
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	logs.Debug("WeChat Event --> ", string(raw))
	if err := xml.Unmarshal(raw, request); err != nil {
		return err
	}
	return nil
}

func (request *Request) checkSignature(req *http.Request) bool {
	ss := sort.StringSlice{
		request.Token,
		req.FormValue("timestamp"),
		req.FormValue("nonce"),
	}
	sort.Strings(ss)          // sort strings by dictionary
	s := strings.Join(ss, "") // concatenate strings
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil)) == req.FormValue("signature")
}
