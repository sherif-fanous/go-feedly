package mapstructure

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	. "github.com/sfanous/go-feedly/pkg/time"
)

func Decode(input interface{}, result interface{}) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: decodeHook,
		Result:     result,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}

	return nil
}

func decodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if to.Kind() == reflect.Struct && to.Name() == "Time" {
		return &Time{
			Time: time.Unix(int64(data.(float64))/1000, 0),
		}, nil
	}

	return data, nil
}
