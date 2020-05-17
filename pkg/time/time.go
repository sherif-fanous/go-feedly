package time

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// Time wraps a time.Time.
type Time struct {
	time.Time
}

// EncodeValues implements the query.Encoder interface.
func (t *Time) EncodeValues(key string, v *url.Values) error {
	v.Set(key, fmt.Sprintf("%d", t.Time.Unix()*1000))

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix() * 1000)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(b []byte) error {
	temp := struct {
		Updated *int64 `json:"updated,omitempty"`
	}{}

	err := json.Unmarshal(b, &temp.Updated)
	if err != nil {
		return err
	}

	if temp.Updated != nil {
		t.Time = time.Unix(*temp.Updated/1000, 0)
	}

	return err
}
