package crcwatch

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type ExtraHeaderStruct struct {
	key    string
	value  string
	origin runtime.ClientRequestWriter
}

func NewExtraHeaderStruct(key, value string, origin runtime.ClientRequestWriter) *ExtraHeaderStruct {
	return &ExtraHeaderStruct{
		key:    key,
		value:  value,
		origin: origin,
	}
}

func (e *ExtraHeaderStruct) WriteToRequest(req runtime.ClientRequest, fmt strfmt.Registry) error {
	if e.origin != nil {
		if err := e.origin.WriteToRequest(req, fmt); err != nil {
			return err
		}
	}
	return req.SetHeaderParam(e.key, e.value)
}

func NewBypassWhiteListHeader(origin runtime.ClientRequestWriter) *ExtraHeaderStruct {
	return NewExtraHeaderStruct("x-bypass-whitelist", "true", origin)
}
