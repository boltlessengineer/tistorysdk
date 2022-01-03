package tistorysdk

import (
	"encoding/json"
	"strconv"
	"time"
)

type SInt int

func (i *SInt) Int() int {
	return int(*i)
}

func (i *SInt) String() string {
	return strconv.Itoa(int(*i))
}

func (i *SInt) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	tmp, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}
	*i = SInt(tmp)
	return nil
}

// SDate is type for date formated as timestamp
type SDate time.Time

func (d *SDate) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	*d = SDate(time.Unix(i, 0))
	return nil
}
