package web

import "encoding/json"

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (r Result[T]) Marshal() []byte {
	json, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return json
}

func (r Result[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
