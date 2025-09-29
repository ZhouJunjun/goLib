/**
 * @author junjunzhou
 * @date 2021/6/8
 */
package t

import (
	"strings"
	"time"
)

type JsonTime time.Time

func (t *JsonTime) Time() time.Time {
	return (time.Time)(*t)
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	ymdHms := YmdHms(time.Time(t).Unix())
	return []byte("\"" + ymdHms + "\""), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	tmpData := strings.Trim(string(data), "\"")
	timestamp := Timestamp(tmpData)
	*t = JsonTime(time.Unix(timestamp, 0))
	return nil
}

func (t JsonTime) String() string {
	return YmdHms(time.Time(t).Unix())
}
