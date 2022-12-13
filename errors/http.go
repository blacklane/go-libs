package errors

import (
	"encoding/json"
	"errors"
	"io"
)

type errorResponse struct {
	Errors []*Error `json:"errors"`
}

func WriteJSON(w io.Writer, err error) {
	res := &errorResponse{
		Errors: innerErrors(err),
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func innerErrors(err error) []*Error {
	list := make([]*Error, 0)
	for err != nil {
		se, ok := err.(*Error)
		if ok {
			list = append(list, se)
		}
		err = errors.Unwrap(err)
	}
	return list
}
