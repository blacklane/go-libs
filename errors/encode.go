package errors

import (
	"encoding/json"
	"errors"
	"io"
)

func WriteJSON(w io.Writer, err error) {
	res := struct {
		Errors []*Error `json:"errors"`
	}{
		Errors: innerErrors(err),
	}
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		panic(err)
	}
}

func innerErrors(err error) []*Error {
	list := make([]*Error, 0)
	for ; err != nil; err = errors.Unwrap(err) {
		if se, ok := err.(*Error); ok {
			list = append(list, se)
		}
	}
	return list
}
