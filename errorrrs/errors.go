// Package errorrs aims to provide more structure and details
// than the standard library package.
package errorrrs

type ID int

const (
	BadRequest ID = iota + 1
	NotFound
	InternalServerError
)

type E struct {
	ID  ID
	Msg string
}
