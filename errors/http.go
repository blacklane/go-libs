package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorResponse struct {
	Errors []*Error `json:"errors"`
}

func Send(w http.ResponseWriter, err error) {
	fe := From(err)

	w.WriteHeader(fe.Code)
	w.Header().Set("content-type", "application/json")

	res := &errorResponse{
		Errors: innerErrors(fe),
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
