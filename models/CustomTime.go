package models

import (
	"encoding/json"
	"time"
)

type CustomTime time.Time

var _ json.Unmarshaler = &CustomTime{}

const dateFormat = "2006-01-02 15:04"

func (mt *CustomTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation(dateFormat, s, time.UTC)
	if err != nil {
		return err
	}
	*mt = CustomTime(t)
	return nil
}

func (mt CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

func (mt *CustomTime) String() string {
	return time.Time(*mt).UTC().String()
}
