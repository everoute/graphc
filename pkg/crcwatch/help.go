package crcwatch

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type ExtraHeaderStruct struct {
	Key   string
	Value string
}

func (e *ExtraHeaderStruct) WriteToRequest(req runtime.ClientRequest, _ strfmt.Registry) error {
	return req.SetHeaderParam(e.Key, e.Value)
}

func NewBypassWhiteListHeader() *ExtraHeaderStruct {
	return &ExtraHeaderStruct{
		Key:   "x-bypass-whitelist",
		Value: "true",
	}
}
