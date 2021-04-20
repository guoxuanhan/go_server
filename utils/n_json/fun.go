package n_json

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Unmarshal(data []byte, i interface{}) error {
	if err := json.Unmarshal(data, i); err != nil {
		fmt.Printf("json unmarshal  %v\n%v\n%v", err, string(data), reflect.TypeOf(i))
		return err
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var data []byte
	var err error

	//if _, ok := v.([]byte); ok {
	//	data = i.([]byte)
	//} else {
	//	if a, ok := v.(string); ok {
	//		data = []byte(a)
	//	} else {
	//		jsonMsg, _ := json.Marshal(v)
	//		data = jsonMsg
	//	}
	//}

	if data, err = json.Marshal(v); err != nil {
		fmt.Printf("json marshal not right  ok ?  %v  :  %v", err, v)
		panic("pppppp")
		return data, err
	}
	return data, err
}
