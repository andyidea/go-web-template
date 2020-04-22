package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	"reflect"
	"regexp"
	"time"
)

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//EncryptPassword 密码MD5加盐加密
func EncryptPassword(password string, salt string) string {
	h := md5.New()
	h.Write([]byte(password))
	if salt != "" {
		h.Write([]byte(salt))
	}
	return hex.EncodeToString(h.Sum(nil))
}

//结构体转map
func StructToMap(obj interface{}) map[string]string {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).String()
	}
	return data
}

func JsonMarshal(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Warn().Err(err).Msg("json marshal failed.")
		return ""
	}

	return string(b)
}

func GetZeroTime(d time.Time) string {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).Format("2006-01-02")
}

func GetEightTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 8, 0, 0, 0, d.Location())
}

func GetYesterdayTime(d time.Time) string {
	return time.Date(d.Year(), d.Month(), d.Day()-1, 0, 0, 0, 0, d.Location()).Format("2006-01-02")
}

func GetTomorrowTime(d time.Time) string {
	return time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, 0, d.Location()).Format("2006-01-02")
}

func GetTomorrowTime1(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, 0, d.Location())
}

func GetNextHour(d time.Time, num int) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour()+num, d.Minute(), d.Second(), 0, d.Location())
}

// 时间转换成日期字符串 (time.Time to "2006-01-02")
func TimeToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 日期字符串转换成时间 ("2006-01-02" to time.Time)
func DateStrToTime(d string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02", d, time.Local)
	return t
}

func DateStr2Time(d string) *time.Time {
	t, _ := time.ParseInLocation("2006-01-02", d, time.Local)
	return &t
}

func TimeStr2Time(d string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", d, time.Local)
	return t
}

func IsContain(items []uint64, item uint64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 获取当日晚上24点（次日0点）的时间
func Get24Time(t time.Time) time.Time {
	dateStr := TimeToDate(t.Add(time.Hour * 24))
	return DateStrToTime(dateStr)
}

func Get14DayLater(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+15, 0, 0, 0, 0, d.Location())
}

func TransSrotOrder(elementOrder string) string {
	if "ascending" == elementOrder {
		return "ASC"
	} else if "descending" == elementOrder {
		return "DESC"
	}
	return ""
}

func StringsIn(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

func UintsIn(us []uint64, target uint64) bool {
	for _, s := range us {
		if s == target {
			return true
		}
	}
	return false
}
func CheckReferrer(c *gin.Context) bool {
	referrer := c.Request.Header.Get("referer")
	matched, _ := regexp.MatchString(`^(https://|http://){1}(localhost|etong-test\.vsattech\.com|etong\.vsattech\.com){1}(/[\w-./?%&=]*)?$`, referrer)
	return matched
}

func GenSonyflake() (uint64, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return id, err
	}
	return id, err
}
