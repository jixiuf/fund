package eastmoney

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type JsonDateTime time.Time

func (jt JsonDateTime) String() string {
	return fmt.Sprintf("\"%s\"", time.Time(jt).Format("2006-01-02"))
}

func (jt JsonDateTime) MarshalJSON() ([]byte, error) {
	if time.Time(jt).IsZero() {
		return []byte(`""`), nil
	}

	var stamp = fmt.Sprintf("\"%s\"", time.Time(jt).Format("2006-01-02"))
	return []byte(stamp), nil
}

func (jt *JsonDateTime) UnmarshalJSON(b []byte) error {
	timeStr := string(b)
	if timeStr == `"0000-00-00T00:00:00"` {
		return nil
	}
	if timeStr == `"0001-01-01T00:00:00"` {
		return nil
	}

	timestr := strings.Trim(timeStr, "\"")
	if timestr == "" || timestr == "0" {
		return nil
	}

	time1, err := time.ParseInLocation("2006-01-02T15:04:05", timestr, time.Local)
	if err != nil {
		fmt.Println(err)
		timeStampI, err := strconv.Atoi(timestr)
		if err != nil {
			return err
		}
		if timeStampI == 0 {
			return nil
		}

		*jt = JsonDateTime(time.Unix(int64(timeStampI), 0))
	} else {
		*jt = JsonDateTime(time1)

	}

	return nil
}

func (jt JsonDateTime) ToTime() time.Time {
	return time.Time(jt)
}
