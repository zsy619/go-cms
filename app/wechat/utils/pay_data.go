package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type SortArray struct {
	k interface{}
	v interface{}
}

type PayData map[string]string

func ParseDate(v map[string]string) PayData {
	return PayData(v)
}

func (data PayData) Set(key, val string) {
	data[key] = val
}

func (data PayData) SetString(key, val string) {
	data.Set(key, val)
}

func (data PayData) SetInt(key string, val int64) {
	data.Set(key, strconv.FormatInt(val, 10))
}

func (data PayData) SetFloat(key string, val float64) {
	data.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
}

func (data PayData) SetBoolean(key string, val bool) {
	data.Set(key, strconv.FormatBool(val))
}

func (data PayData) SetDate(key string, val time.Time) {
	loc, _ := time.LoadLocation(DATE_TIMEZONE)
	data.Set(key, val.In(loc).Format(DATE_TIME_FORMAT))
}

func (data PayData) Get(key string) string {
	return data[key]
}

func (data PayData) IsExist(key string) bool {
	_, b := data[key]
	return b
}

func (data PayData) SortKeys() []string {
	var keys sort.StringSlice
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (data PayData) ToJson() string {
	b, e := json.Marshal(&data)
	if e != nil {
		return ""
	}
	return string(b)
}

func (data PayData) ToMap() map[string]string {
	return data
}

func (data PayData) IsNil() bool {
	return len(data) == 0
}

// CurrentTimeStampMS get current time with millisecond
func CurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

// CurrentTimeStampNS get current time with nanoseconds
func CurrentTimeStampNS() int64 {
	return time.Now().UnixNano()
}

// CurrentTimeStamp get current time with unix
func CurrentTimeStamp() int64 {
	return time.Now().Unix()
}

func CurrentTimeStampString() string {
	return strconv.FormatInt(CurrentTimeStamp(), 10)
}

func SHA1(s string) string {
	m := sha1.New()
	m.Write([]byte(s))
	return fmt.Sprintf("%x", m.Sum(nil))
}
