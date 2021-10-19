// Package errorrs aims to provide more structure and details
// than the standard library package.
package errorrrs

import (
	"context"
	"encoding/json"
	"net/http"
)

type ID int

const (
	BadRequest ID = iota + 1
	NotFound
	InternalServerError
)

var _ error = (*E)(nil)

type E struct {
	ID  ID     `json:"-"`
	Msg string `json:"error"`
}

func (e *E) Error() string {
	return e.Msg
}

// GokitErrorEncoder is an implementation of Gokit func signature
// for writing (and "minor" handling) HTTP error responses
func GokitErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	switch errt := err.(type) {
	case *E:
		var hs int
		if errt.ID == BadRequest {
			hs = http.StatusBadRequest
		} else if errt.ID == NotFound {
			hs = http.StatusNotFound
		} else {
			hs = http.StatusInternalServerError
		}
		w.WriteHeader(hs)
		bits, err := json.Marshal(errt)
		if err != nil {
			panic(err.Error())
		}
		w.Write(bits)
	default:
		e := E{Msg: errt.Error()}
		bits, err := json.Marshal(e)
		if err != nil {
			panic(err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bits)
	}
}
