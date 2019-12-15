package handler

import (
	"github.com/golang/protobuf/ptypes/any"
)

type Event interface {
	UUID() string
	Any() *any.Any
}

type event struct {
	uuid string
	any  *any.Any
}

func (e event) UUID() string {
	return e.uuid
}

func (e event) Any() *any.Any {
	return e.any
}
